use sha2::{Digest, Sha256};

pub fn hash_string(inp: String) -> String {
    let mut h = Sha256::new();
    h.update(&inp.as_bytes());
    hex::encode(h.finalize())
}
