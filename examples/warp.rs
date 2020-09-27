use warp::Filter;

#[tokio::main]
async fn main() {
    let hello = warp::path!("hello" / String)
        .map(|name| format!("Hello, {}!", name));
    let health = warp::path!(".within" / "health")
        .map(|| "OK");
    let routes = hello.or(health);

    warp::serve(routes)
        .run(([0, 0, 0, 0], 3030))
        .await;
}
