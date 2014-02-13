use std::run::{ProcessOptions};
use std::path::Path;

pub struct Config {
  root_dir: Path,
  work_dir: Path,
  pid_file: Option<~str>,
  uid: u32,
  mem: uint,
  cpu: uint,
  command: ~[~str],
}

#[cfg(target_os = "linux")]
#[cfg(target_os = "android")]
pub fn run(c: Config) {
  println!("Linux");
}

#[cfg(target_os = "freebsd")]
#[cfg(target_os = "macos")]
pub fn run(c: Config) {
  println!("Dumb jail")
  println!("root: {:?}", c.root_dir.as_str().unwrap());
  println!("work: {:?}", c.work_dir.as_str().unwrap());
  println!("pid_file: {}", c.pid_file);
  println!("uid: {}", c.uid);
  println!("mem: {}", c.mem);
  println!("cpu: {}", c.cpu);
  println!("command: {:?}", c.command);

  let p = ProcessOptions{
    env: None, // TODO: Change the PATH and stuff
    dir: Some(&c.work_dir),
    in_fd: Some(0), // TODO: close stdin
    out_fd: Some(1),
    err_fd: Some(2),
  };
  // TODO: See how to create a new process group
}
