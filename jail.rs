pub use jail::common::Config;
pub use jail::common::look_path;
pub use jail::macos::run;

mod common {
  use std::path::Path;
  use std::path::BytesContainer;

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
      let paths = PATHS.map(|p| join_path(p.as_bytes(), cmd.clone()));
      match paths.iter().find(|p| join_path(*root.clone(), (**p).container_as_bytes()).is_file()) {
        Some(m) => Some(m.clone()),
        None => None
      }
    }
  }

  // Like path.join but doesn't replace by p2 if it's absolute
  pub fn join_path<T1: BytesContainer, T2: BytesContainer>(p1: T1, p2: T2) -> Path {
    let mut p = Path::new(p1).join("foo");
    p.set_filename(p2);
    p
  }
}

#[cfg(target_os = "macos")]
pub mod macos {
  use libc::uint64_t;
  use libc::c_char;
  use libc::c_int;
  use libc::c_void;
  use jail::common::Config;
  use jail::common::join_path;
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
    let program = match c.program.is_absolute() {
      true => join_path(c.root_dir.clone(), c.program.clone()),
      false => join_path(c.work_dir.clone(), c.program.clone()),
    };

    println!("Dumb jail")
    println!("root: {:?}", c.root_dir.as_str().unwrap());
    println!("work: {:?}", c.work_dir.as_str().unwrap());
    println!("uid: {}", c.uid);
    println!("mem: {}", c.mem);
    println!("cpu: {}", c.cpu);
    println!("net: {}", c.net);
    println!("program: {:?}", program.as_str().unwrap());

    

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
