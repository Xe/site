use super::{PostNotFound, SeriesNotFound};
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
