const curry = (fn, ...presetArgs) => {
  return (...args) => {
    const mergedArgs = presetArgs.concat(args);

    if (mergedArgs.length >= fn.length) {
      return fn(...mergedArgs);
    }

    return curry(fn, ...mergedArgs);
  };
};

const compose = (...fns) => (input) => {
  return fns.reduce((output, fn) => fn(output), input);
};

const throttle = (fn, timeout = 1000) => {
  let timer;
  return (...args) => {
    clearTimeout(timer);
    timer = setTimeout(() => fn(...args), timeout);
  };
};

export {
  curry,
  compose,
  throttle,
};
