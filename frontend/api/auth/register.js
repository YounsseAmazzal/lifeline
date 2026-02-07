 //  Load Cities (Mockup for now)
    document.addEventListener("DOMContentLoaded", () => {
        const citySelect = document.getElementById('citySelect');
        const cities = ["الدار البيضاء", "الرباط", "مراكش", "فاس", "طنجة"];
        citySelect.innerHTML = '<option value="" disabled selected>اختر المدينة...</option>';
        cities.forEach(c => {
            let opt = document.createElement('option');
            opt.value = c; opt.innerText = c; citySelect.appendChild(opt);
        });
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

        const registerData = {
            name: `${fname} ${lname}`,      
            userName: email,                 
            email: email,
            phoneNumber: phone,
            password: pass,
            bloodGroup: blood,
            city: city,
            country: "Morocco"              
        };

        try {
            const response = await auth.register(registerData);
            
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