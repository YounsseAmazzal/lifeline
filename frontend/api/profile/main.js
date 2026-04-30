function formatDateInput(dateString) {
    if (!dateString) return "";
    const date = new Date(dateString);
    if (Number.isNaN(date.getTime())) return "";
    return date.toISOString().slice(0, 10);
}

function applyLanguageFromProfile(user) {
    if (!user || !user.language) return;
    if (typeof persistLanguage === "function") persistLanguage(user.language);
    if (typeof applyLanguageLayout === "function") applyLanguageLayout(user.language);
}

function updateAvailabilityUI(available) {
    const ping = document.getElementById('availPing');
    const dot = document.getElementById('availDot');
    const text = document.getElementById('availText');
    const badge = document.getElementById('availBadge');
    if (!ping) return;

    if (available) {
        ping.className = 'animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75';
        dot.className = 'relative inline-flex rounded-full h-2.5 w-2.5 bg-green-500';
        text.textContent = 'Available to donate';
        text.className = 'text-xs font-bold text-green-700';
        badge.className = 'inline-flex items-center gap-2.5 bg-green-50 border border-green-100 px-3 py-1.5 rounded-full';
    } else {
        ping.className = 'hidden';
        dot.className = 'relative inline-flex rounded-full h-2.5 w-2.5 bg-slate-400';
        text.textContent = 'Not available';
        text.className = 'text-xs font-bold text-slate-500';
        badge.className = 'inline-flex items-center gap-2.5 bg-slate-100 border border-slate-200 px-3 py-1.5 rounded-full';
    }
}

document.addEventListener("DOMContentLoaded", async () => {
    try {
        const token = localStorage.getItem('lifeline_token');
        if (!token) {
            window.location.href = '../auth/login.html';
            return;
        }

        const user = await accountApi.profile();
        applyLanguageFromProfile(user);

        const nameEl   = document.getElementById('userName');
        const cityEl   = document.getElementById('userCity');
        const bloodEl  = document.getElementById('userBlood');
        const avatarEl = document.getElementById('userAvatar');
        const toggle   = document.getElementById('availabilityToggle');

        if (nameEl)   nameEl.innerText  = user.name       || "—";
        if (cityEl)   cityEl.innerText  = user.city       || "—";
        if (bloodEl)  bloodEl.innerText = user.blood_group || "--";
        if (avatarEl && user.photo_url) avatarEl.src = user.photo_url;
        if (toggle) toggle.checked = !!user.available;
        updateAvailabilityUI(!!user.available);

        // Keep identity card in sync when toggle changes
        if (toggle) {
            toggle.addEventListener('change', () => updateAvailabilityUI(toggle.checked));
        }

        const form = document.getElementById('profileForm');
        if (form) {
            form.elements.name.value         = user.name         || "";
            form.elements.email.value        = user.email        || "";
            form.elements.phone_number.value = user.phone_number || "";
            form.elements.area.value         = user.area         || "";
            form.elements.state.value        = user.state        || "";
            form.elements.country.value      = user.country      || "Morocco";
            form.elements.postal_code.value  = user.postal_code  || "";
            form.elements.gender.value       = user.gender       || "";
            form.elements.blood_group.value  = user.blood_group  || "";
            form.elements.date_of_birth.value = formatDateInput(user.date_of_birth);
            form.elements.available.checked  = !!user.available;
            form.elements.language.value     = user.language || (typeof getPreferredLanguage === "function" ? getPreferredLanguage() : "ar");

            await loadCities(user.city || "");
        }
    } catch (error) {
        if (error.message === "Unauthorized") {
            localStorage.removeItem('lifeline_token');
            window.location.href = '../auth/login.html';
            return;
        }
        console.error(error);
    }
});

async function loadCities(currentCity) {
    const citySelect = document.getElementById('profileCitySelect');
    if (!citySelect) return;
    citySelect.innerHTML = '<option value="">Loading...</option>';
    try {
        const cities = await geoApi.cities();
        citySelect.innerHTML = '<option value="">Select city...</option>';
        cities.forEach((city) => {
            const option = document.createElement('option');
            option.value = city.name;
            option.textContent = city.name;
            citySelect.appendChild(option);
        });
        if (currentCity) citySelect.value = currentCity;
    } catch (_) {
        citySelect.innerHTML = '<option value="">Failed to load cities</option>';
    }
}

async function updateProfile(event) {
    event.preventDefault();

    const form      = document.getElementById('profileForm');
    const submitBtn = document.getElementById('saveProfileBtn');
    const originalHTML = submitBtn.innerHTML;
    submitBtn.disabled = true;
    submitBtn.innerHTML = '<i class="fa-solid fa-circle-notch fa-spin mr-2"></i> Saving...';

    try {
        const formData = new FormData();
        formData.append('name',          form.elements.name.value);
        formData.append('email',         form.elements.email.value);
        formData.append('phone_number',  form.elements.phone_number.value);
        formData.append('city',          form.elements.city.value);
        formData.append('area',          form.elements.area.value);
        formData.append('state',         form.elements.state.value);
        formData.append('country',       form.elements.country.value);
        formData.append('postal_code',   form.elements.postal_code.value);
        formData.append('gender',        form.elements.gender.value);
        formData.append('blood_group',   form.elements.blood_group.value);
        formData.append('date_of_birth', form.elements.date_of_birth.value);
        formData.append('available',     form.elements.available.checked ? "true" : "false");
        formData.append('language',      form.elements.language.value || (typeof getPreferredLanguage === "function" ? getPreferredLanguage() : "ar"));

        if (navigator.geolocation) {
            await new Promise((resolve) => {
                navigator.geolocation.getCurrentPosition(
                    (pos) => {
                        formData.append('latitude',  String(pos.coords.latitude));
                        formData.append('longitude', String(pos.coords.longitude));
                        resolve();
                    },
                    () => resolve(),
                    { enableHighAccuracy: false, maximumAge: 60000, timeout: 5000 }
                );
            });
        }

        const photoInput = document.getElementById('newProfilePhoto');
        if (photoInput && photoInput.files && photoInput.files[0]) {
            formData.append('photo', photoInput.files[0]);
        }

        const updated = await accountApi.updateProfile(formData);
        applyLanguageFromProfile(updated);
        window.showAutoAlert('Profile updated successfully', 'success');
        setTimeout(() => window.location.reload(), 1200);
    } catch (error) {
        window.showAutoAlert(error.message || 'Failed to update profile', 'error');
    } finally {
        submitBtn.disabled = false;
        submitBtn.innerHTML = originalHTML;
    }
}

function logout() {
    localStorage.removeItem('lifeline_token');
    window.location.href = '../auth/login.html';
}
