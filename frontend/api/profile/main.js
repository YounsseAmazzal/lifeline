
document.addEventListener("DOMContentLoaded", async () => {
    try {
        const token = localStorage.getItem('lifeline_token');
        if (!token) {
            window.location.href = 'login.html';
            return;
        }
        const user = await apiRequest('/account/profile', 'GET');
                const els = {
            name: document.getElementById('userName'),
            city: document.getElementById('userCity'),
            blood: document.getElementById('userBlood'),
            avatar: document.getElementById('userAvatar'),
            toggle: document.getElementById('availabilityToggle')
        };

        if (els.name) els.name.innerText = user.name;
        if (els.city) els.city.innerText = user.city;
        if (els.blood) els.blood.innerText = user.blood_group;
        // console.log(user.blood_group)
        // console.log(user.date_of_birth)
        if (els.avatar && user.photoUrl) {
            els.avatar.src = user.photoUrl;
        }

        if (els.toggle) els.toggle.checked = user.available;

    } catch (error) {
            if (error.message === "Unauthorized") {
        localStorage.removeItem('lifeline_token');
        window.location.href = 'login.html';
    }
    }
});

function logout() {
    localStorage.removeItem('lifeline_token');
    window.location.href = '../auth/login.html';
}