(function () {
  const LANG_KEY = 'lifeline_lang';

  function normalizeLanguage(lang) {
    const value = String(lang || '').toLowerCase();
    if (value === 'ar' || value === 'en' || value === 'fr') {
      return value;
    }
    return 'ar';
  }

  function getPreferredLanguage() {
    const saved = localStorage.getItem(LANG_KEY);
    if (saved) {
      return normalizeLanguage(saved);
    }

    const browser = (navigator.language || 'ar').slice(0, 2);
    return normalizeLanguage(browser);
  }

  function persistLanguage(lang) {
    localStorage.setItem(LANG_KEY, normalizeLanguage(lang));
  }

  function applyLanguageLayout(lang) {
    const value = normalizeLanguage(lang);
    const html = document.documentElement;
    const body = document.body;

    html.setAttribute('lang', value);
    html.setAttribute('dir', value === 'ar' ? 'rtl' : 'ltr');

    if (!body) return;

    if (value === 'ar') {
      body.classList.remove('font-sans');
      body.classList.add('font-arabic');
    } else {
      body.classList.remove('font-arabic');
      body.classList.add('font-sans');
    }
  }

  window.getPreferredLanguage = getPreferredLanguage;
  window.persistLanguage = persistLanguage;
  window.applyLanguageLayout = applyLanguageLayout;

  document.addEventListener('DOMContentLoaded', function () {
    applyLanguageLayout(getPreferredLanguage());
  });
})();
