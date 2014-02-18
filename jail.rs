pub use jail::common::Config;
pub use jail::common::look_path;
pub use jail::macos::run;

mod common {
  use std::path::Path;

  pub struct Config {
    root_dir: Path,
    work_dir: Path,
    uid: u32,
    mem: uint,
    cpu: uint,
    net: bool,
    program: Path,
    args: ~[~str],
  }

  static PATHS: [&'static str, .. 4] = ["/bin", "/sbin", "/usr/bin", "/usr/sbin"];

  // Tries to find an executable
  // TODO: Look for the executable bit for the select user/group
  pub fn look_path(program: ~str, root: ~Path) -> Option<Path> {
    let cmd = Path::new(program.clone());
    if cmd.is_absolute() || program.contains("/") { // TODO: Use std::path::SEP
      Some(cmd)
    } else {
      let paths = PATHS.map(|p| Path::new(p.as_bytes()).join(cmd.clone()));
      match paths.iter().find(|p| root.join(*p).is_file()) {
        Some(m) => Some(m.clone()),
        None => None
      }
    }
  }
}

#[cfg(target_os = "macos")]
pub mod macos {
  use libc::uint64_t;
  use libc::c_char;
  use libc::c_int;
  use libc::c_void;
  use jail::common::Config;
  use std::io::process::ProcessConfig;

  //static kSBXProfileNoInternet = "no-internet";
  static kSBXProfileNoNetwork: &'static str = "no-network";
  //static kSBXProfileNoWrite = "no-write";
  //static kSBXProfileNoWriteExceptTemporary
  //static kSBXProfilePureComputation

  static SANDBOX_NAMED: uint64_t = 1;

  extern {
    fn sandbox_init(profile: *c_char, flags: uint64_t, errorbuf: *mut *mut c_char) -> c_int;
    fn sandbox_free_error(errorbuf: *c_char) -> c_void;
  }

  pub fn run(c: ~Config) {
    println!("Dumb jail")
    println!("root: {:?}", c.root_dir.as_str().unwrap());
    println!("work: {:?}", c.work_dir.as_str().unwrap());
    println!("uid: {}", c.uid);
    println!("mem: {}", c.mem);
    println!("cpu: {}", c.cpu);
    println!("net: {}", c.net);
    println!("program: {:?}", c.program);

    let program = match c.program.is_absolute() {
      true => c.root_dir.join(c.program.clone()),
      false => c.work_dir.join(c.program.clone()),
    };

    // if !c.net {
    //   kSBXProfileNoNetwork.to_c_str().with_ref(|profile| {
    //     unsafe {
    //       let errorbuf: *mut *mut c_char = (); // FIXME: null pointers FTW
    //       if sandbox_init(profile, 0, errorbuf) } < 0 {
    //         fail!(*errorbuf)
    //       }
    //     }
    //   });
    // }

    // let p = ProcessConfig{ 
    //   program: program.as_str().unwrap(), 
    //   args: c.args, 
    //   env: None, // TODO: Change the PATH and stuff 
    //   cwd: c.work_dir.as_str(), 
    //   io: [], 
    // }; 

    
    // TODO: See how to create a new process group
  }
}
