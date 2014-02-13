extern mod getopts;

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

fn main() {
  let args = os::args();
  let program = args[0].clone();
  let opts = ~[
    optopt("r",  "root", "Root directory (r/o) to run the command under", "DIR"),
    optopt("w",  "work", "The work directory and PWD (r/w)", "DIR"),
    optopt("u",  "user", "Run the command under. Disabled on suid exec", "NAME"),
    optopt("m",  "mem",  "Max memory. Unlimited (0) by default", "BYTES"),
    optopt("c",  "cpu",  "How many cpu shares to allocate. Unlimited (0) by default.", "NUM"),
    optopt("p",  "pid",  "Write the command's pid there if specified", "PATH"),
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

  let mut c = jail::Config {
    work_dir: Path::new(~"/tmp"),
    root_dir: Path::new(~"/"),
    pid_file: None,
    uid: 0,
    mem: 0,
    cpu: 0,
    command: ~[],
  };

  match matches.opt_str("root") {
    Some(m) => {
      c.root_dir = Path::new(m)
    }
    None => {
      fail_usage(program, opts, "root is missing");
      return
    }
  };
  if ! c.root_dir.is_dir() {
    fail_usage(program, opts, "root is not a directory");
  }

  match matches.opt_str("work") {
    Some(m) => {
      c.work_dir = Path::new(m)
    }
    None => {
      // TODO: Make a tmp dir instead
      fail_usage(program, opts, "work is missing");
      return;
    }
  };
  if ! c.work_dir.is_dir() {
    fail_usage(program, opts, "work is not a directory");
    return
  }

  if matches.free.is_empty() {
    fail_usage(program, opts, "command is missing");
    return;
  }

  c.uid = unsafe { libc::getuid() };
  if c.uid == 0 {
    match matches.opt_str("user") {
      Some(m) => {
        // TODO: resolve the uid for that user
        println!("user: {}", m);
      }
      None => () // keep root as the runner
    }
  } else {
    match matches.opt_str("user") {
      Some(_) => println!("ignoring user option"),
      None => ()
    }
  }

  c.mem = match matches.opt_str("mem") {
    Some(m) => from_str::<uint>(m).unwrap(),
    None => 0
  };

  c.cpu = match matches.opt_str("cpu") {
    Some(m) => from_str::<uint>(m).unwrap(),
    None => 0
  };

  c.pid_file = matches.opt_str("pid");
  c.command = matches.free;

    /* do_work(root_dir, work_dir); */
  jail::run(c)
}
