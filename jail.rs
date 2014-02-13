
pub struct Config {
  root_dir: ~str,
  work_dir: ~str,
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
  println!("root: {}", c.root_dir);
  println!("work: {}", c.work_dir);
  println!("pid_file: {}", c.pid_file);
  println!("uid: {}", c.uid);
  println!("mem: {}", c.mem);
  println!("cpu: {}", c.cpu);
  println!("command: {}", c.command.connect(", "));
}
