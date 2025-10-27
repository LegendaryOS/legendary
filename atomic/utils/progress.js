const chalk = require('chalk');

let direction = 1;
let length = 1;
let position = 0;
const maxLength = 20;
const wall = '|';

exports.startProgress = () => {
  process.stdout.write(chalk.cyan('Progress: '));
  const interval = setInterval(() => {
    process.stdout.clearLine(0);
    process.stdout.cursorTo(0);
    process.stdout.write(chalk.cyan('Progress: '));

    let bar = '';
    for (let i = 0; i < length; i++) {
      if (i === position) {
        bar += direction > 0 ? '>' : '<';
      } else {
        bar += '-';
      }
    }
    process.stdout.write(wall + bar.padEnd(maxLength, ' ') + wall);

    position += direction;
    if (position >= length - 1 || position <= 0) {
      direction *= -1;
      if (position >= length - 1 && length < maxLength) {
        length++;
      }
    }
  }, 100);
  return interval;
};

exports.stopProgress = (interval) => {
  if (interval) {
    clearInterval(interval);
    process.stdout.clearLine(0);
    process.stdout.cursorTo(0);
    process.stdout.write('\n');
  }
};
