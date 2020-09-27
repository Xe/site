use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, Deserialize, Serialize)]
pub struct Author {
    pub id: i32,
    pub name: String,
}

#[derive(Clone, Debug, Deserialize, Serialize)]
pub struct Comment {
    pub id: i32,
    pub author: Author,
    pub body: String,
    pub in_reply_to: i32,
}

fn main() {
    let data = r#"
  {
    "id": 31337,
    "author": {
      "id": 420,
      "name": "Cadey"
    },
    "body": "hahaha its is an laughter image",
    "in_reply_to": 31335
  }
  "#;

    let c: Comment = serde_json::from_str(data).expect("json to parse");
    println!("comment: {:#?}", c);
}
