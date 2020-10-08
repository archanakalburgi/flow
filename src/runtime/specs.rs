use crate::catalog;
use estuary_protocol::{consumer, protocol};
use std::collections::BTreeMap;
use std::fmt::Write;

#[derive(Debug)]
pub struct DerivationSet(BTreeMap<String, ()>);

impl std::convert::TryFrom<consumer::ListResponse> for DerivationSet {
    type Error = ();

    // Eventually we'll want to hoist existing shards into a DerivationSet
    // that knows about current Etcd revisions and splits.
    // For now, we assume a single shard and over-write it each time.
    fn try_from(_: consumer::ListResponse) -> Result<Self, Self::Error> {
        Ok(DerivationSet(BTreeMap::new()))
    }
}

impl DerivationSet {
    pub fn update_from_catalog(&mut self, db: &catalog::DB) -> catalog::Result<()> {
        let mut stmt =
            db.prepare("SELECT collection_name FROM collection_details WHERE is_derivation")?;
        let mut rows = stmt.query(rusqlite::NO_PARAMS)?;

        while let Some(row) = rows.next()? {
            self.0.insert(row.get(0)?, ());
        }
        Ok(())
    }

    pub fn build_recovery_log_apply_request(&self) -> protocol::ApplyRequest {
        let changes = self
            .0
            .iter()
            .map(|(collection, _)| {
                let labels = Some(protocol::LabelSet {
                    labels: [
                        ("app.gazette.dev/managed-by", "estuary.dev/flow"),
                        ("content-type", "application/x-gazette-recoverylog"),
                    ]
                    .iter()
                    .map(|(n, v)| protocol::Label {
                        name: (*n).to_owned(),
                        value: (*v).to_owned(),
                    })
                    .collect(),
                });

                let fragment = Some(protocol::journal_spec::Fragment {
                    length: 1 << 28, // 256MB.
                    compression_codec: (protocol::CompressionCodec::None as i32),
                    stores: vec!["file:///".to_owned()],
                    refresh_interval: Some(std::time::Duration::from_secs(5 * 60).into()),
                    retention: None,
                    flush_interval: None,
                    path_postfix_template: String::new(),
                });

                protocol::apply_request::Change {
                    upsert: Some(protocol::JournalSpec {
                        name: format!("recovery/derivation/{}/00/00000000", collection),
                        replication: 1,
                        labels,
                        fragment,
                        flags: 0,
                        max_append_rate: 0,
                    }),
                    expect_mod_revision: -1, // TODO (always update).
                    delete: String::new(),
                }
            })
            .collect::<Vec<_>>();

        protocol::ApplyRequest { changes }
    }

    pub fn build_shard_apply_request(&self, catalog_url: &str) -> consumer::ApplyRequest {
        let changes = self
            .0
            .iter()
            .map(|(collection, _)| {
                let labels = Some(protocol::LabelSet {
                    labels: [
                        ("app.gazette.dev/managed-by", "estuary.dev/flow"),
                        ("estuary.dev/catalog-url", catalog_url),
                        ("estuary.dev/derivation", collection.as_ref()),
                        ("estuary.dev/key-begin", "00"),
                        ("estuary.dev/key-end", "ffffffff"),
                        ("estuary.dev/rclock-begin", "0000000000000000"),
                        ("estuary.dev/rclock-end", "ffffffffffffffff"),
                    ]
                    .iter()
                    .map(|(n, v)| protocol::Label {
                        name: (*n).to_owned(),
                        value: (*v).to_owned(),
                    })
                    .collect(),
                });

                consumer::apply_request::Change {
                    upsert: Some(consumer::ShardSpec {
                        id: format!("{}/00/00000000", collection),
                        sources: Vec::new(),
                        recovery_log_prefix: "recovery/derivation".to_owned(),
                        hint_prefix: "/estuary/flow/hints".to_owned(),
                        hint_backups: 2,
                        max_txn_duration: Some(prost_types::Duration {
                            seconds: 1,
                            nanos: 0,
                        }),
                        min_txn_duration: None,
                        disable: false,
                        hot_standbys: 0, // TODO
                        disable_wait_for_ack: false,
                        labels,
                    }),
                    expect_mod_revision: -1, // TODO (always update).
                    delete: String::new(),
                }
            })
            .collect::<Vec<_>>();

        consumer::ApplyRequest {
            changes,
            ..Default::default()
        }
    }
}

fn _hex_key(key: &[u8]) -> String {
    let mut s = String::with_capacity(2 * key.len());
    for byte in key {
        write!(s, "{:02X}", byte).unwrap();
    }
    s
}

fn _hex_rc(rc: u32) -> String {
    format!("{:08x}", rc)
}