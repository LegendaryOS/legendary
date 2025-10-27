const { execSync } = require('child_process');

exports.createSnapshot = (snapRoot) => {
  execSync(`btrfs subvolume snapshot / ${snapRoot}`);
};

exports.prepareChroot = (snapRoot) => {
  ['proc', 'sys', 'dev', 'run'].forEach(dir => {
    execSync(`mkdir -p ${snapRoot}/${dir}`);
    execSync(`mount --bind /${dir} ${snapRoot}/${dir}`);
  });
  // Bind /etc/resolv.conf for network
  execSync(`mount --bind /etc/resolv.conf ${snapRoot}/etc/resolv.conf`);
};

exports.runInChroot = (snapRoot, cmd) => {
  execSync(`chroot ${snapRoot} /bin/bash -c "${cmd}"`);
};

exports.cleanupChroot = (snapRoot) => {
  ['proc', 'sys', 'dev', 'run'].forEach(dir => {
    execSync(`umount ${snapRoot}/${dir}`);
  });
  execSync(`umount ${snapRoot}/etc/resolv.conf`);
};

exports.deploySnapshot = (snapRoot) => {
  const id = execSync(`btrfs subvolume show ${snapRoot} | grep "Subvolume ID" | awk '{print $3}'`).toString().trim();
  execSync(`btrfs subvolume set-default ${id} /`);
  // Note: Reboot required to apply
};
