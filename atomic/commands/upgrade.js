const { execSync } = require('child_process');
const chalk = require('chalk');

module.exports = () => {
  console.log(chalk.green('Running LegendaryOS upgrade script'));
  try {
    execSync('/usr/share/LegendaryOS/SCRIPTS/legendaryos-upgrade');
    console.log(chalk.green('Upgrade script executed successfully.'));
  } catch (error) {
    console.error(chalk.red(`Error during upgrade: ${error.message}`));
  }
};
