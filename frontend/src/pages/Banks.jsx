import React, { useState, useEffect } from 'react';

const Banks = () => {
    const [banks, setBanks] = useState([]);
    const [loading, setLoading] = useState(true);

    // Todo: Move to a separate API service
    useEffect(() => {
        const fetchBanks = async () => {
            try {
                const response = await fetch('/api/banks', {
                    headers: {
                        // Include token if endpoint is protected (optional for banks usually, but good practice)
                        'Authorization': `Bearer ${localStorage.getItem('token')}`
                    }
                });

                if (!response.ok) {
                    throw new Error('Failed to fetch banks');
                }

                const data = await response.json();
                // Ensure data is an array
                setBanks(Array.isArray(data) ? data : []);
            } catch (err) {
                console.error(err);
            } finally {
                setLoading(false);
            }
        }

        fetchBanks();
    }, []);

    if (loading) {
        return <div className="container" style={{ paddingTop: '4rem', textAlign: 'center' }}>Loading...</div>;
    }

    return (
        <div className="container" style={{ padding: '4rem 0' }}>
            <h2 style={{ marginBottom: '2rem' }}>Blood Banks Near You</h2>

            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(350px, 1fr))', gap: '2rem' }}>
                {banks.map(bank => (
                    <div key={bank.id} className="card">
                        <h3 style={{ fontSize: '1.25rem', marginBottom: '0.5rem' }}>{bank.name}</h3>
                        <p style={{ color: 'var(--text-secondary)', marginBottom: '1rem' }}>
                            📍 {bank.area}, {bank.city}
                        </p>

                        <div style={{ marginBottom: '1.5rem' }}>
                            <div style={{ display: 'flex', flexWrap: 'wrap', gap: '0.5rem' }}>
                                {bank.blood_groups.map((bg, idx) => (
                                    <span key={idx} style={{
                                        padding: '0.25rem 0.5rem',
                                        borderRadius: '4px',
                                        fontSize: '0.85rem',
                                        backgroundColor: bg.quantity === 0 ? '#FEF2F2' : '#EFF6FF',
                                        color: bg.quantity === 0 ? '#DC2626' : '#2563EB',
                                        border: `1px solid ${bg.quantity === 0 ? '#FECACA' : '#BFDBFE'}`
                                    }}>
                                        {bg.group}: {bg.quantity > 0 ? 'Available' : 'Critical'}
                                    </span>
                                ))}
                            </div>
                        </div>

                        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                            <a href={`tel:${bank.phone_number}`} className="btn btn-secondary" style={{ padding: '0.5rem 1rem', fontSize: '0.9rem' }}>
                                Call {bank.phone_number}
                            </a>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default Banks;
