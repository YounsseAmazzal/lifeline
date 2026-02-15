const API_URL = "http://localhost:8080/api"; 

 async function apiRequest(endpoint, method, data = null) {
    const headers = { "Content-Type": "application/json" };

    const token = localStorage.getItem("lifeline_token");
    if (token) {
        headers["Authorization"] = `Bearer ${token}`;
    }

    const config = {
        method,
        headers,
    };

    if (data) config.body = JSON.stringify(data);

    const response = await fetch(`${API_URL}${endpoint}`, config);

    if (response.status === 401) {
        localStorage.removeItem('lifeline_token');
        throw new Error("Unauthorized");
    }

    let result = null;
    const contentType = response.headers.get("content-type");

    if (contentType && contentType.includes("application/json")) {
        result = await response.json();
    }

    if (!response.ok) {
        throw new Error(result?.error || "خطأ في الاتصال");
    }

    return result;
}

 async function apiRequestForm(endpoint, method, formData) {
    const headers = {};
    const token = localStorage.getItem("lifeline_token");
    if (token) {
        headers["Authorization"] = `Bearer ${token}`;
    }

    const response = await fetch(`${API_URL}${endpoint}`, {
        method,
        headers,
        body: formData
    });

    if (response.status === 401) {
        localStorage.removeItem('lifeline_token');
        throw new Error("Unauthorized");
    }

    let result = null;
    const contentType = response.headers.get("content-type");
    if (contentType && contentType.includes("application/json")) {
        result = await response.json();
    }

    if (!response.ok) {
        throw new Error(result?.error || "خطأ في الاتصال");
    }

    return result;
}

const auth = {
    login: (userName, password) => apiRequest("/account/login", "POST", { userName, password }),
    register: (userData) => userData instanceof FormData
        ? apiRequestForm("/account/register", "POST", userData)
        : apiRequest("/account/register", "POST", userData)
};

const accountApi = {
    profile: () => apiRequest("/account/profile", "GET"),
    updateProfile: (formData) => apiRequestForm("/account/profile", "PUT", formData),
    uploadPhoto: (formData) => apiRequestForm("/account/profile/photo", "POST", formData)
};

const geoApi = {
    cities: () => apiRequest("/geo/cities", "GET")
};
console.log(geoApi)
