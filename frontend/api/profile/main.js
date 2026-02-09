function formatDateInput(dateString) {
    if (!dateString) return "";
    const date = new Date(dateString);
    if (Number.isNaN(date.getTime())) return "";
    return date.toISOString().slice(0, 10);
}

function applyLanguageFromProfile(user) {
    if (!user || !user.language) return;
    if (typeof persistLanguage === "function") {
        persistLanguage(user.language);
    }
    if (typeof applyLanguageLayout === "function") {
        applyLanguageLayout(user.language);
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

        const els = {
            name: document.getElementById('userName'),
            city: document.getElementById('userCity'),
            blood: document.getElementById('userBlood'),
            avatar: document.getElementById('userAvatar'),
            toggle: document.getElementById('availabilityToggle')
        };

        if (els.name) els.name.innerText = user.name || "-";
        if (els.city) els.city.innerText = user.city || "-";
        if (els.blood) els.blood.innerText = user.blood_group || "--";
        if (els.avatar && user.photo_url) {
            els.avatar.src = user.photo_url;
        }
        if (els.toggle) els.toggle.checked = !!user.available;

        const form = document.getElementById('profileForm');
        if (form) {
            form.elements.name.value = user.name || "";
            form.elements.email.value = user.email || "";
            form.elements.phone_number.value = user.phone_number || "";
            form.elements.city.value = user.city || "";
            form.elements.area.value = user.area || "";
            form.elements.state.value = user.state || "";
            form.elements.country.value = user.country || "Morocco";
            form.elements.postal_code.value = user.postal_code || "";
            form.elements.gender.value = user.gender || "";
            form.elements.blood_group.value = user.blood_group || "";
            form.elements.date_of_birth.value = formatDateInput(user.date_of_birth);
            form.elements.available.checked = !!user.available;

            form.elements.language.value = user.language || (typeof getPreferredLanguage === "function" ? getPreferredLanguage() : "ar");

            await loadCities(form.elements.city.value);
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

    const form = document.getElementById('profileForm');
    const submitBtn = document.getElementById('saveProfileBtn');
    const submitLabel = submitBtn.innerHTML;
    submitBtn.disabled = true;
    submitBtn.innerHTML = 'Saving...';

    try {
        const formData = new FormData();
        const boolAvailable = form.elements.available.checked ? "true" : "false";

        formData.append('name', form.elements.name.value);
        formData.append('email', form.elements.email.value);
        formData.append('phone_number', form.elements.phone_number.value);
        formData.append('city', form.elements.city.value);
        formData.append('area', form.elements.area.value);
        formData.append('state', form.elements.state.value);
        formData.append('country', form.elements.country.value);
        formData.append('postal_code', form.elements.postal_code.value);
        formData.append('gender', form.elements.gender.value);
        formData.append('blood_group', form.elements.blood_group.value);
        formData.append('date_of_birth', form.elements.date_of_birth.value);
        formData.append('available', boolAvailable);
        formData.append('language', form.elements.language.value || (typeof getPreferredLanguage === "function" ? getPreferredLanguage() : "ar"));

        if (navigator.geolocation) {
            await new Promise((resolve) => {
                navigator.geolocation.getCurrentPosition(
                    (pos) => {
                        formData.append('latitude', String(pos.coords.latitude));
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
        alert('Profile updated successfully');
        window.location.reload();
    } catch (error) {
        alert(error.message || 'Failed to update profile');
    } finally {
        submitBtn.disabled = false;
        submitBtn.innerHTML = submitLabel;
    }
}

function logout() {
    localStorage.removeItem('lifeline_token');
    window.location.href = '../auth/login.html';
}
