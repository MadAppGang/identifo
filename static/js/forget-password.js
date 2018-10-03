const emailRegex = /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/

const makeObserver = () => {
  const listeners = []
  const subscribe = (listener) => {
    const index = listeners.length;
    listeners.push(listener);
    let isUnsubscribed = false;

    const unsubscribe = () => {
      if (isUnsubscribed) return;

      listeners.splice(idnex, 1);
      isUnsubscribed = true;
    }

    return unsubscribe;
  }

  const emit = (...args) => listeners.forEach(cb => cb(...args));
  return { subscribe, emit };
}

const getErrorMessage = (err) => {
  if (!err || !err.message) {
    return 'error';
  }

  if (!err.response || !err.response.data) {
    return err.message;
  }

  return err.response.data.error;
};

const createFormView = (formId, emailFieldId, buttonId, errorMessageId) => {
  const form = document.getElementById(formId);
  const emailField = document.getElementById(emailFieldId);
  const button = document.getElementById(buttonId);
  const errorMessage = document.getElementById(errorMessageId);
  const { subscribe, emit } = makeObserver();

  if (!form || !emailField || !button || !errorMessage) {
    throw new Error('Element not found');
  }

  const displayError = (text) => {
    errorMessage.innerHTML = text;
    errorMessage.style.opacity = 1;
  }

  hideError = () => errorMessage.style.opacity = 0;

  const sendRequest = (data) => {
    return axios.post('http://127.0.0.1:8080/password/forgot', data);
  }

  const validateEmail = (value) => {
    return emailRegex.test(value);
  }

  const handleClickOnSubmit = (e) => {
    e.preventDefault();
    const { value } = emailField;
    const isValid = validateEmail(value);

    if (!isValid) {
      displayError('*Invalid Email');
      return;
    }

    hideError();

    sendRequest({ username: value.toLowerCase() })
      .then(() => emit())
      .catch(err => displayError(getErrorMessage(err)));
  }

  button.addEventListener('click', handleClickOnSubmit);
  return {
    rootEl: form,
    onSuccess: subscribe,
  }
}

const createFinalView = (elemId) => {
  const elem = document.getElementById(elemId);

  if (!elem) {
    throw new Error('Element not found');
  }

  return {
    rootEl: elem,
  };
}

const createScreen = (view) => {
  const hiddenClass = 'hidden';
  const hide = () => view.rootEl.classList.add(hiddenClass);
  const show = () => view.rootEl.classList.remove(hiddenClass);
  return {
    hide, show, view
  };
}

const init = () => {
  const finalView = createFinalView('final');
  const formView = createFormView('form', 'email', 'submit', 'error');

  const finalScreen = createScreen(finalView);
  const formScreen = createScreen(formView);
  
  finalScreen.hide();

  formScreen.view.onSuccess(() => {
    formScreen.hide();
    finalScreen.show();
  });
};

init();
