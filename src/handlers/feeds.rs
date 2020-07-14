use crate::{
    app::State,
};
use lazy_static::lazy_static;
use prometheus::{IntCounterVec, register_int_counter_vec, opts};
use std::{sync::Arc};
use warp::{
    Reply,
};

lazy_static! {
    static ref HIT_COUNTER: IntCounterVec =
        register_int_counter_vec!(opts!("feed_hits", "Number of hits to various feeds"), &["kind"])
        .unwrap();
}

pub fn jsonfeed(state: Arc<State>) -> impl Reply {
    HIT_COUNTER.with_label_values(&["json"]).inc();
    let state = state.clone();
    warp::reply::json(&state.jf)
}
