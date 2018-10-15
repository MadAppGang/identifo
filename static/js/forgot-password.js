"use strict";

var emailRegex = /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
var makeObserver = function makeObserver() {
  var listeners = [];
  var subscribe = function subscribe(listener) {
    var index = listeners.length;
    listeners.push(listener);
    var isUnsubscribed = false;
    var unsubscribe = function unsubscribe() {
      if (isUnsubscribed) return;
      listeners.splice(idnex, 1);
      isUnsubscribed = true;
    };
    return unsubscribe;
  };
  var emit = function emit() {
    for (
      var _len = arguments.length, args = Array(_len), _key = 0;
      _key < _len;
      _key++
    ) {
      args[_key] = arguments[_key];
    }

    return listeners.forEach(function(cb) {
      return cb.apply(undefined, args);
    });
  };
  return { subscribe: subscribe, emit: emit };
};
var getErrorMessage = function getErrorMessage(err) {
  if (!err || !err.message) {
    return "error";
  }
  if (!err.response || !err.response.data) {
    return err.message;
  }
  return err.response.data.error;
};
var createFormView = function createFormView(
  formId,
  emailFieldId,
  buttonId,
  errorMessageId
) {
  var form = document.getElementById(formId);
  var emailField = document.getElementById(emailFieldId);
  var button = document.getElementById(buttonId);
  var errorMessage = document.getElementById(errorMessageId);

  var _makeObserver = makeObserver(),
    subscribe = _makeObserver.subscribe,
    emit = _makeObserver.emit;

  if (!form || !emailField || !button || !errorMessage) {
    throw new Error("Element not found");
  }
  var displayError = function displayError(text) {
    errorMessage.innerHTML = text;
    errorMessage.style.opacity = 1;
  };
  var hideError = function hideError() {
    return (errorMessage.style.opacity = 0);
  };
  var sendRequest = function sendRequest(data) {
    return axios.post("http://127.0.0.1:8080/password/forgot", data);
  };
  var validateEmail = function validateEmail(value) {
    return emailRegex.test(value);
  };
  var handleClickOnSubmit = function handleClickOnSubmit(e) {
    e.preventDefault();
    var value = emailField.value;

    var isValid = validateEmail(value);
    if (!isValid) {
      displayError("*Invalid Email");
      return;
    }
    hideError();
    sendRequest({ username: value.toLowerCase() })
      .then(function() {
        return emit();
      })
      .catch(function(err) {
        return displayError(getErrorMessage(err));
      });
  };
  button.addEventListener("click", handleClickOnSubmit);
  return {
    rootEl: form,
    onSuccess: subscribe
  };
};
var createFinalView = function createFinalView(elemId) {
  var elem = document.getElementById(elemId);
  if (!elem) {
    throw new Error("Element not found");
  }
  return {
    rootEl: elem
  };
};
var createScreen = function createScreen(view) {
  var hiddenClass = "hidden";
  var hide = function hide() {
    return view.rootEl.classList.add(hiddenClass);
  };
  var show = function show() {
    return view.rootEl.classList.remove(hiddenClass);
  };
  return {
    hide: hide,
    show: show,
    view: view
  };
};
var init = function init() {
  var finalView = createFinalView("final");
  var formView = createFormView("form", "email", "submit", "error");
  var finalScreen = createScreen(finalView);
  var formScreen = createScreen(formView);

  finalScreen.hide();
  formScreen.view.onSuccess(function() {
    formScreen.hide();
    finalScreen.show();
  });
};
init();
