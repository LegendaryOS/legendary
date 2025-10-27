const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');
const chalk = require('chalk');
const { createSnapshot, prepareChroot, runInChroot, cleanupChroot, deploySnapshot } = require('../utils/btrfs');
const { startProgress, stopProgress } = require('../utils/progress');

module.exports = () => {
  console.log(chalk.green('Updating all packages'));

  const timestamp = new Date().toISOString().replace(/[:.-]/g, '');
  const snapDir = `/var/lib/legendary/snapshots/${timestamp}`;
  const snapRoot = path.join(snapDir, '@');

  let progressInterval;
  try {
    fs.mkdirSync(snapDir, { recursive: true });
    createSnapshot(snapRoot);
    prepareChroot(snapRoot);

    // Start progress bar
    progressInterval = startProgress();

    runInChroot(snapRoot, `pacman -Syu --noconfirm`);

    stopProgress(progressInterval);
    cleanupChroot(snapRoot);
    deploySnapshot(snapRoot);
    console.log(chalk.green('Update complete. Reboot to apply changes.'));
    // Optionally clean old snapshots
    require('./clean')();
  } catch (error) {
    stopProgress(progressInterval);
    console.error(chalk.red(`Error during update: ${error.message}`));
    execSync(`btrfs subvolume delete ${snapRoot}`);
    fs.rmdirSync(snapDir);
  }
};
