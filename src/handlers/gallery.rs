use super::PostNotFound;
use crate::{
    app::State,
    post::Post,
    templates::{self, Html, RenderRucte},
};
use std::sync::Arc;
use warp::{http::Response, Rejection, Reply};

pub async fn index(state: Arc<State>) -> Result<impl Reply, Rejection> {
    let state = state.clone();
    Response::builder().html(|o| templates::galleryindex_html(o, state.gallery.clone()))
}

pub async fn post_view(name: String, state: Arc<State>) -> Result<impl Reply, Rejection> {
    let mut want: Option<Post> = None;

    for post in &state.gallery {
        if post.link == format!("gallery/{}", name) {
            want = Some(post.clone());
        }
    }

    match want {
        None => Err(PostNotFound("blog".into(), name).into()),
        Some(post) => {
            let body = Html(post.body_html.clone());
            Response::builder().html(|o| templates::gallerypost_html(o, post, body))
        }
    }
}
