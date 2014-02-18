pub use jail::macos::run;
pub use jail::common::Config;

mod common {
  use std::path::Path;
  pub struct Config {
    root_dir: Path,
    work_dir: Path,
    pid_file: Option<~str>,
    uid: u32,
    mem: uint,
    cpu: uint,
    net: bool,
    command: ~[~str],
  }
}

#[cfg(target_os = "macos")]
pub mod macos {
  use libc::uint64_t;
  use libc::c_char;
  use libc::c_int;
  use libc::c_void;
  use libc::uintptr_t;
  use jail::common::Config;
  use std::io::process::ProcessConfig;

  //static kSBXProfileNoInternet = "no-internet";
  static kSBXProfileNoNetwork: &'static str = "no-network";
  //static kSBXProfileNoWrite = "no-write";
  //static kSBXProfileNoWriteExceptTemporary
  //static kSBXProfilePureComputation

  static SANDBOX_NAMED: uint64_t = 1;

  static PATHS: [&'static str, .. 4] = ["bin", "sbin", "usr/bin", "usr/sbin"];

  extern {
    // FIXME: **char type for errorbuf
    fn sandbox_init(profile: *c_char, flags: uint64_t, errorbuf: uintptr_t) -> c_int;
    fn sandbox_free_error(errorbuf: *c_char) -> c_void;
  }

  pub fn prepare(c: &Config) {
    if !c.net {
      kSBXProfileNoNetwork.to_c_str().with_ref(|profile| {
        unsafe { sandbox_init(profile, 0, 0 as uintptr_t) }
      });
    }
  }

  pub fn run(c: ~Config) {
    println!("Dumb jail")
    println!("root: {:?}", c.root_dir.as_str().unwrap());
    println!("work: {:?}", c.work_dir.as_str().unwrap());
    println!("pid_file: {}", c.pid_file);
    println!("uid: {}", c.uid);
    println!("mem: {}", c.mem);
    println!("cpu: {}", c.cpu);
    println!("command: {:?}", c.command);

    prepare(c);

    /* let p = ProcessConfig{ */
    /*   program: "TODO", */
    /*   args: [], */
    /*   env: None, // TODO: Change the PATH and stuff */
    /*   cwd: c.work_dir.as_str(), */
    /*   io: [], */
    /* }; */

    
    // TODO: See how to create a new process group
  }
}
