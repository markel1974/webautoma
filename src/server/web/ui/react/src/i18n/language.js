import en from './languages/en';
//import it from './languages/it';

import jQuery from 'jquery';
require('jquery.cookie');

const languages = {
  'en': en,
  //'it': it
};

const supported = [
  { name: 'English', code: 'en' }
  //,{ name: 'Italiano', code: 'it' }
];

let language;
let locale = jQuery.cookie('locale') || navigator.language || 'en';
//locale = locale.indexOf('en') === -1 ? 'it' : 'en';
locale = 'en';
setLocale(locale);

function setLocale(loc) {
  const loc2 = loc.toLowerCase();
  if (languages[loc2]) {
    locale = loc;
    language = languages[loc2]
  }
}

function getLocale() {
  return locale
}

/**
 * @return {string}
 */
function L(k) {
  let t = language[k];
  if (typeof t === 'undefined') {
    t = k;
  }

  let tr = '';
  let inTag = false;
  let num = '';
  let stash = '';
  for (let i = 0; i < t.length; i++) {
    //const cc = t[i].charCodeAt();
    const cc = t.charCodeAt(i);
    if (cc === 123) {
      if (inTag) {
        tr += stash;
        stash = '';
      }
      inTag = true;
    } else if (inTag && cc >= 48 && cc <= 57) {
      num += t[i];
    } else if (inTag && cc === 125) {
      const index = parseInt(num);
      if (arguments.length - 1 > index) {
        tr += arguments[index + 1];
      }
      num = '';
      stash = '';
      inTag = false;
      continue;
    } else {
      inTag = false;
      tr += stash;
      stash = '';
    }

    if (inTag) {
      stash += t[i];
    } else {
      tr += t[i];
    }
  }

  return tr;
}

const lang = { getLocale, setLocale, supported, L };
export default lang;
