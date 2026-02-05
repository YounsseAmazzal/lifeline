import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';

const Register = () => {
    const navigate = useNavigate();
    const [formData, setFormData] = useState({
        name: '',
        username: '',
        email: '',
        phone_number: '',
        password: '',
        date_of_birth: '',
        gender: 'male',
        blood_group: 'O+',
        area: '',
        city: '',
        state: '',
        country: 'Morocco',
        postal_code: ''
    });

    const handleChange = (e) => {
        setFormData({
            ...formData,
            [e.target.name]: e.target.value
        });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();

        // Format Date to ISO strict
        const formattedData = {
            ...formData,
            date_of_birth: new Date(formData.date_of_birth).toISOString(),
        };

        try {
            const response = await fetch('/api/account/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(formattedData),
            });

            const data = await response.json();

            if (!response.ok) {
                alert(data.error || 'Registration failed');
                return;
            }

            // Success
            alert('Registration successful! Please login.');
            navigate('/login');
        } catch (err) {
            console.error(err);
            alert('Network error occurred.');
        }
    };

    return (
        <div className="container" style={{ padding: '4rem 0' }}>
            <div className="card" style={{ maxWidth: '800px', margin: '0 auto' }}>
                <h2 style={{ marginBottom: '2rem', textAlign: 'center' }}>Create an Account</h2>

                <form onSubmit={handleSubmit} style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1.5rem' }}>

                    {/* Section: Personal Info */}
                    <div style={{ gridColumn: 'span 2' }}>
                        <h3 style={{ fontSize: '1.1rem', color: 'var(--primary)', marginBottom: '1rem', borderBottom: '1px solid #eee', paddingBottom: '0.5rem' }}>Personal Information</h3>
                    </div>

                    <div className="input-group">
                        <label className="input-label">Full Name</label>
                        <input type="text" name="name" className="input-field" required onChange={handleChange} />
                    </div>

                    <div className="input-group">
                        <label className="input-label">Date of Birth</label>
                        <input type="date" name="date_of_birth" className="input-field" required onChange={handleChange} />
                    </div>

                    <div className="input-group">
                        <label className="input-label">Gender</label>
                        <select name="gender" className="input-field" onChange={handleChange}>
                            <option value="male">Male</option>
                            <option value="female">Female</option>
                        </select>
                    </div>

                    <div className="input-group">
                        <label className="input-label">Blood Group</label>
                        <select name="blood_group" className="input-field" onChange={handleChange}>
                            {['A+', 'A-', 'B+', 'B-', 'AB+', 'AB-', 'O+', 'O-'].map(bg => (
                                <option key={bg} value={bg}>{bg}</option>
                            ))}
                        </select>
                    </div>

                    {/* Section: Contact & Address */}
                    <div style={{ gridColumn: 'span 2', marginTop: '1rem' }}>
                        <h3 style={{ fontSize: '1.1rem', color: 'var(--primary)', marginBottom: '1rem', borderBottom: '1px solid #eee', paddingBottom: '0.5rem' }}>Contact & Location</h3>
                    </div>

                    <div className="input-group">
                        <label className="input-label">Phone Number</label>
                        <input type="tel" name="phone_number" className="input-field" required onChange={handleChange} />
                    </div>

                    <div className="input-group">
                        <label className="input-label">Email (Optional)</label>
                        <input type="email" name="email" className="input-field" onChange={handleChange} />
                    </div>

                    <div className="input-group">
                        <label className="input-label">City</label>
                        <input type="text" name="city" className="input-field" required onChange={handleChange} />
                    </div>

                    <div className="input-group">
                        <label className="input-label">Area / District</label>
                        <input type="text" name="area" className="input-field" required onChange={handleChange} />
                    </div>

                    <div className="input-group">
                        <label className="input-label">State / Region</label>
                        <input type="text" name="state" className="input-field" required onChange={handleChange} />
                    </div>

                    <div className="input-group">
                        <label className="input-label">Postal Code</label>
                        <input type="text" name="postal_code" className="input-field" required onChange={handleChange} />
                    </div>

                    {/* Section: Security */}
                    <div style={{ gridColumn: 'span 2', marginTop: '1rem' }}>
                        <h3 style={{ fontSize: '1.1rem', color: 'var(--primary)', marginBottom: '1rem', borderBottom: '1px solid #eee', paddingBottom: '0.5rem' }}>Security</h3>
                    </div>

                    <div className="input-group">
                        <label className="input-label">Username</label>
                        <input type="text" name="username" className="input-field" required onChange={handleChange} />
                    </div>

                    <div className="input-group">
                        <label className="input-label">Password</label>
                        <input type="password" name="password" className="input-field" required minLength="4" onChange={handleChange} />
                    </div>

                    <div style={{ gridColumn: 'span 2', marginTop: '1rem' }}>
                        <button type="submit" className="btn btn-primary" style={{ width: '100%' }}>Create Account</button>
                    </div>

                </form>

                <p style={{ marginTop: '1.5rem', textAlign: 'center', fontSize: '0.9rem', color: 'var(--text-secondary)' }}>
                    Already have an account? <Link to="/login" style={{ color: 'var(--primary)', fontWeight: '600' }}>Login instead</Link>
                </p>
            </div>
        </div>
    );
};

export default Register;
