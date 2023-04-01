const std = @import("std");
const Builder = std.build.Builder;
const CrossTarget = std.zig.CrossTarget;

pub fn build(b: *std.build.Builder) void {
    // Standard release options allow the person running `zig build` to select
    // between Debug, ReleaseSafe, ReleaseFast, and ReleaseSmall.
    const mode = b.standardReleaseOptions();
    const target = CrossTarget{ .cpu_arch = .wasm32, .os_tag = .wasi };

    const exe = b.addExecutable("hello", "src/hello.zig");
    exe.setTarget(target);
    exe.setBuildMode(mode);
    exe.install();
}
