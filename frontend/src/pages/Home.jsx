import React from 'react';
import { Link } from 'react-router-dom';

const Home = () => {
    return (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '4rem', paddingBottom: '4rem' }}>

            {/* Hero Section */}
            <section style={{
                backgroundColor: '#FFF1F2', // Very light red
                padding: '6rem 0',
                textAlign: 'center'
            }}>
                <div className="container">
                    <h1 style={{
                        fontSize: '3.5rem',
                        fontWeight: '800',
                        marginBottom: '1.5rem',
                        color: 'var(--secondary)',
                        lineHeight: 1.1
                    }}>
                        The "Uber" for <br /> <span style={{ color: 'var(--primary)' }}>Blood Donation</span>
                    </h1>
                    <p style={{
                        fontSize: '1.25rem',
                        color: 'var(--text-secondary)',
                        maxWidth: '600px',
                        margin: '0 auto 2.5rem auto'
                    }}>
                        Connect instantly with donors and blood banks.
                        Smart location matching to save lives when it matters most.
                    </p>
                    <div style={{ display: 'flex', gap: '1rem', justifyContent: 'center' }}>
                        <Link to="/register" className="btn btn-primary" style={{ padding: '1rem 2rem', fontSize: '1.1rem' }}>
                            Find Blood Now
                        </Link>
                        <Link to="/register" className="btn btn-secondary" style={{ padding: '1rem 2rem', fontSize: '1.1rem' }}>
                            Become a Donor
                        </Link>
                    </div>
                </div>
            </section>

            {/* Stats / Info */}
            <section className="container">
                <div style={{
                    display: 'grid',
                    gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))',
                    gap: '2rem'
                }}>
                    <div className="card">
                        <h3 style={{ fontSize: '1.5rem', marginBottom: '1rem', color: 'var(--secondary)' }}>📍 Smart Location</h3>
                        <p style={{ color: 'var(--text-secondary)' }}>
                            We use PostGIS technology to find donors within a specific 5km radius of the hospital. No more spamming entire cities.
                        </p>
                    </div>
                    <div className="card">
                        <h3 style={{ fontSize: '1.5rem', marginBottom: '1rem', color: 'var(--secondary)' }}>🛡️ Privacy First</h3>
                        <p style={{ color: 'var(--text-secondary)' }}>
                            Your phone number is hidden. Communication happens securely through the platform to prevent harassment.
                        </p>
                    </div>
                    <div className="card">
                        <h3 style={{ fontSize: '1.5rem', marginBottom: '1rem', color: 'var(--secondary)' }}>⚡ Instant Alerts</h3>
                        <p style={{ color: 'var(--text-secondary)' }}>
                            Donors receive push notifications only when their specific blood type is needed nearby.
                        </p>
                    </div>
                </div>
            </section>

        </div>
    );
};

export default Home;
