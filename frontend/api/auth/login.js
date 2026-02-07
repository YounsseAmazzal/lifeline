document.getElementById('loginForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        
        // UI Feedback
        const btn = e.target.querySelector('button');
        const originalText = btn.innerHTML;
        btn.innerHTML = '<i class="fa-solid fa-circle-notch fa-spin"></i> جاري الاتصال...';
        btn.disabled = true;

        // Get Values
        const email = document.getElementById('email').value; 
        const password = document.getElementById('password').value;

 try {
    const response = await auth.login(email, password);
    
    localStorage.setItem("lifeline_token", response.token);
    
    if (response.role === "Admin") {
        window.location.href = "admin.html";
    } else if (response.role === "Sponsor") {
        window.location.href = "sponsor.html";
    } else {
        window.location.href = "dashboard.html";
    }

} catch (error) {
    alert("Error: " + error.message);
            btn.innerHTML = originalText;
            btn.disabled = false;
        }
    });