#[no_uv];
extern crate native;
extern crate getopts;

use getopts::{optopt,optflag,getopts,usage,OptGroup};
use std::os;
use std::libc;
use std::path::Path;

mod jail;

/* fn do_work(inp: &str, out: Option<~str>) { */
/*     println!(inp); */
/*     println(match out { */
/*         Some(x) => x, */
/*         None => ~"No Output" */
/*     }); */
/* } */

fn print_usage(program: &str, opts: &[OptGroup]) {
  println!("{}", usage(format!("Usage: {} [options] -- command ...", program), opts));
}

fn fail_usage(program: &str, opts: &[OptGroup], error: &str) {
  println!("{}\n", error);
  print_usage(program, opts);
  os::set_exit_status(2);
}

#[start]
fn start(argc: int, argv: **u8) -> int { native::start(argc, argv, main) }

fn main() {
  let args = os::args();
  let program = args[0].clone();
  let opts = ~[
    optopt("r",  "root", "Root directory (r/o) to run the command under", "DIR"),
    optopt("w",  "work", "The work directory and PWD (r/w)", "DIR"),
    optopt("u",  "user", "Run the command under. Disabled on suid exec", "NAME"),
    optopt("m",  "mem",  "Max memory. Unlimited (0) by default", "BYTES"),
    optopt("c",  "cpu",  "How many cpu shares to allocate. Unlimited (0) by default.", "NUM"),
    optflag("h", "help", "Prints this help"),
  ];

  let matches = match getopts(args.tail(), opts) {
    Ok(m) => { m }
    Err(f) => {
      fail_usage(program, opts, f.to_err_msg());
      return;
    }
  };

  if matches.opt_present("help") {
    print_usage(program, opts);
    return;
  }

  let root_dir = match matches.opt_str("root") {
    Some(m) => Path::new(m),
    None => Path::new(~"/")
  };
  if ! root_dir.is_dir() {
    fail_usage(program, opts, "root is not a directory");
    return;
  }

  let work_dir = match matches.opt_str("work") {
    Some(m) => Path::new(m),
    None => os::getcwd()
  };
  if ! work_dir.is_dir() {
    fail_usage(program, opts, "work is not a directory");
    return
  }

  if matches.free.is_empty() {
    fail_usage(program, opts, "command is missing");
    return;
  }

  let uid = unsafe { libc::getuid() };
  match matches.opt_str("user") {
    Some(m) => {
      // Only change uid if we're root. suid only allows current user.
      if uid == 0 {
        // TODO: resolve the uid for that user
        println!("user: {}", m);
      } else {
        println!("ignoring user option");
      }
    }
    None => ()
  }

  let mem = match matches.opt_str("mem") {
    Some(m) => from_str::<uint>(m).unwrap(),
    None => 0
  };

  let cpu = match matches.opt_str("cpu") {
    Some(m) => from_str::<uint>(m).unwrap(),
    None => 0
  };

  let args = matches.free;

  let program = match jail::look_path(args[0].clone(), ~root_dir.clone()) {
    Some(m) => m,
    None => {
      fail_usage(program, opts, "program not found");
      return
    }
  };

  let c = ~jail::Config {
    root_dir: root_dir,
    work_dir: work_dir,
    uid: 0,
    mem: 0,
    cpu: 0,
    net: true,
    program: program,
    args: args,
  };

    /* do_work(root_dir, work_dir); */
  jail::run(c)
}
