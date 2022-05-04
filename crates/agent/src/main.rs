use anyhow::Context;
use clap::Parser;
use futures::{FutureExt, TryFutureExt};
use tracing::info;

/// Agent is a daemon which runs server-side tasks of the Flow control-plane.
#[derive(Parser, Debug)]
#[clap(author, version, about, long_about = None)]
struct Args {
    /// URL of the postgres database.
    #[clap(
        long = "database",
        env = "DATABASE_URL",
        default_value = "postgres://flow:flow@127.0.0.1:5432/control_development"
    )]
    database_url: url::Url,
    /// Path to CA certificate of the database.
    #[clap(long = "database-ca", env = "DATABASE_CA")]
    database_ca: Option<String>,
    /// URL of the builds root.
    #[clap(
        long = "builds-root",
        env = "BUILDS_ROOT",
        default_value = "file:///var/tmp/"
    )]
    builds_root: url::Url,
    /// URL of the data-plane Gazette broker.
    #[clap(
        long = "broker-address",
        env = "BROKER_ADDRESS",
        default_value = "http://localhost:8080"
    )]
    _broker_address: url::Url,
    /// URL of the data-plane Flow consumer.
    #[clap(
        long = "consumer-address",
        env = "CONSUMER_ADDRESS",
        default_value = "http://localhost:9000"
    )]
    _consumer_address: url::Url,
    /// Docker network for connector invocations.
    #[clap(long = "connector-network", default_value = "host")]
    connector_network: String,
    /// Path to binaries like `flowctl`.
    #[clap(long = "bin-dir")]
    bindir: String,
}

#[tokio::main]
async fn main() -> Result<(), anyhow::Error> {
    let args = Args::parse();
    info!(?args, "started!");

    // Use reasonable defaults for printing structured logs to stderr.
    let subscriber = tracing_subscriber::FmtSubscriber::builder()
        .with_env_filter(tracing_subscriber::EnvFilter::from_default_env())
        .finish();
    tracing::subscriber::set_global_default(subscriber).expect("setting tracing default failed");

    let mut pg_options = args
        .database_url
        .as_str()
        .parse::<sqlx::postgres::PgConnectOptions>()
        .context("parsing database URL")?
        .application_name("agent");

    // If a database CA was provided, require that we use TLS with full cert verification.
    if let Some(ca) = &args.database_ca {
        pg_options = pg_options
            .ssl_mode(sqlx::postgres::PgSslMode::VerifyFull)
            .ssl_root_cert(ca);
    } else {
        // Otherwise, prefer TLS but don't require it.
        pg_options = pg_options.ssl_mode(sqlx::postgres::PgSslMode::Prefer);
    }

    let pg_pool = sqlx::postgres::PgPool::connect_with(pg_options)
        .await
        .context("connecting to database")?;

    // Start a logs sink into which agent loops may stream logs.
    let (logs_tx, logs_rx) = tokio::sync::mpsc::channel(8192);
    let logs_sink = agent::logs::serve_sink(pg_pool.clone(), logs_rx);

    let serve_fut = agent::serve(
        vec![
            Box::new(agent::PublishHandler::new(
                &args.connector_network,
                &args.bindir,
                &logs_tx,
                &args.builds_root,
            )),
            Box::new(agent::TagHandler::new(
                &args.connector_network,
                &args.bindir,
                &logs_tx,
            )),
            Box::new(agent::DiscoverHandler::new(
                &args.connector_network,
                &args.bindir,
                &logs_tx,
            )),
        ],
        pg_pool.clone(),
        tokio::signal::ctrl_c().map(|_| ()),
    );

    std::mem::drop(logs_tx);
    let ((), ()) = tokio::try_join!(serve_fut, logs_sink.map_err(Into::into))?;

    Ok(())
}
