use anyhow::Context;
use tokio::signal::unix;

mod capture;
mod inspect;
mod materialize;
mod rpc;

#[derive(clap::Parser, Debug)]
#[clap(about = "Command to start connector proxies for Flow runtime.")]
pub struct Args {
    /// The path (in the container) to the JSON file that contains the
    /// inspection results from the connector image.
    /// Normally produced via command "docker inspect <image>".
    #[clap(short, long)]
    pub image_inspect_json_path: String,

    /// Port on which to listen for requests from the runtime.
    #[clap(short, long, default_value = "8080")]
    pub port: u16,
}

pub async fn run(args: Args) -> anyhow::Result<()> {
    let image = inspect::Image::parse_from_json_file(&args.image_inspect_json_path)
        .context("reading image inspect JSON")?;
    let mut entrypoint = image.get_argv()?;

    // TODO(johnny): Remove this in preference of always using the LOG_LEVEL variable.
    if let Ok(log_level) = std::env::var("LOG_LEVEL") {
        entrypoint.push(format!("--log.level={log_level}"));
    }

    let capture = proto_grpc::capture::driver_server::DriverServer::new(capture::Driver {
        entrypoint: entrypoint.clone(),
    });
    let materialize =
        proto_grpc::materialize::driver_server::DriverServer::new(materialize::Driver {
            entrypoint: entrypoint.clone(),
        });

    let addr = format!("0.0.0.0:{}", args.port).parse().unwrap();

    // Gracefully exit on either SIGINT (ctrl-c) or SIGTERM.
    let mut sigint = unix::signal(unix::SignalKind::interrupt()).unwrap();
    let mut sigterm = unix::signal(unix::SignalKind::terminate()).unwrap();

    let signal = async move {
        tokio::select! {
            _ = sigint.recv() => (),
            _ = sigterm.recv() => (),
        }
        tracing::info!("caught signal to exit");
    };

    let () = tonic::transport::Server::builder()
        .add_service(capture)
        .add_service(materialize)
        .serve_with_shutdown(addr, signal)
        .await?;

    Ok(())
}
