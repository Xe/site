use crate::templates::{self, RenderRucte};
use warp::{
    http::{Response, StatusCode},
    path, Filter, Rejection, Reply,
};

pub async fn index() -> Result<impl Reply, Rejection> {
    Response::builder().html(|o| templates::index_html(o))
}
