// Load cities from backend database
document.addEventListener("DOMContentLoaded", async () => {
    const citySelect = document.getElementById('citySelect');
    citySelect.innerHTML = '<option value="" disabled selected>Loading...</option>';

    try {
        const cities = await geoApi.cities();
        // console.log(cities[0])
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

    //  Handle Registration
    document.getElementById('registerForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const btn = e.target.querySelector('button');
        const originalText = btn.innerHTML;
        btn.innerHTML = '<i class="fa-solid fa-circle-notch fa-spin"></i> Loading...';
        btn.disabled = true;

        const inputs = document.querySelectorAll('input'); 
        const fname = inputs[0].value; // First Name
        const lname = inputs[1].value; // Last Name
        const email = inputs[2].value; // Email
        const phone = inputs[3].value; // Phone
        const pass  = inputs[4].value; // Password
        
        const selects = document.querySelectorAll('select');
        const blood = selects[0].value; 
        const city  = selects[1].value; 

        const formData = new FormData();
        formData.append("name", `${fname} ${lname}`);
        formData.append("userName", email);
        formData.append("email", email);
        formData.append("phoneNumber", phone);
        formData.append("password", pass);
        formData.append("bloodGroup", blood);
        formData.append("city", city);
        formData.append("country", "Morocco");
        if (typeof getPreferredLanguage === "function") {
            formData.append("language", getPreferredLanguage());
        }

        const photoInput = document.getElementById("profilePhoto");
        if (photoInput && photoInput.files && photoInput.files[0]) {
            formData.append("photo", photoInput.files[0]);
        }

        try {
            const response = await auth.register(formData);
            if (response.language && typeof persistLanguage === "function") {
                persistLanguage(response.language);
            }
            
            // Success
            localStorage.setItem("lifeline_token", response.token);
            alert("Mabrouk! Compte t-sayeb.");
            window.location.href = 'login.html';

        } catch (error) {
            alert("Mochkil: " + error.message);
            btn.innerHTML = originalText;
            btn.disabled = false;
        }
    });
