// Load cities from backend database
// import { showAutoAlert } from "../global.js";
document.addEventListener("DOMContentLoaded", async () => {
    const citySelect = document.getElementById('citySelect');
    citySelect.innerHTML = '<option value="" disabled selected>Loading...</option>';

    try {
        const cities = await geoApi.cities();
        // console.log(cities)
        // for (let i=0;i<4;i++){
        //     console.log(cities[i])
        // }
        citySelect.innerHTML = '<option value="" disabled selected>Select city...</option>';
        cities.forEach((c) => {
            const opt = document.createElement('option');
            opt.value = c.ville;
        // console.log(c.name)
            opt.innerText = c.ville;
            citySelect.appendChild(opt);
        });
    } catch (error) {
        citySelect.innerHTML = '<option value="" disabled selected>Failed to load cities</option>';
    }
});

// Handle Registration
document.getElementById('registerForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const btn = e.target.querySelector('button');
    const originalText = btn.innerHTML;
    const citySelect = document.getElementById('citySelect');
    const bloodSelect = document.getElementById('bloodSelect');
    
    if (!citySelect.value || !bloodSelect.value) {
         window.showAutoAlert("Marjo ikhtiyar l'mdina wa fasilat dam", "error");
         return;
    }

    btn.innerHTML = '<i class="fa-solid fa-circle-notch fa-spin"></i> Loading...';
    btn.disabled = true;

    const inputs = document.querySelectorAll('input'); 
    const fname = inputs[0].value; 
    const lname = inputs[1].value; 
    const email = inputs[2].value; 
    const phone = inputs[3].value; 
    const pass  = inputs[4].value;
    
    const formData = new FormData();
    formData.append("name", `${fname} ${lname}`);
    formData.append("userName", email);
    formData.append("email", email);
    formData.append("phoneNumber", phone);
    formData.append("password", pass);
    formData.append("bloodGroup", bloodSelect.value);
    formData.append("city", citySelect.value);
    formData.append("country", "Morocco");

    // Language handling
    if (typeof getPreferredLanguage === "function") {
        formData.append("language", getPreferredLanguage());
    }

    // Photo handling
    const photoInput = document.getElementById("profilePhoto");
    if (photoInput && photoInput.files && photoInput.files[0]) {
        formData.append("photo", photoInput.files[0]);
        console.log("---------------"+photoInput.files[0])
    }
    // console.log(photoInput)
    try {
        const response = await auth.register(formData);
        
        if (response.language && typeof persistLanguage === "function") {
            persistLanguage(response.language);
        }
        
        localStorage.setItem("lifeline_token", response.token);
        window.showAutoAlert("Mabrouk! Compte t-sayeb.", "success");
        setTimeout(() => {
            window.location.href = 'login.html';
        }, 2000);

    } catch (error) {
         window.showAutoAlert("Mochkil: " + error.message, "error"); 
        
        btn.innerHTML = originalText;
        btn.disabled = false;
    }
});
