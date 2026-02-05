import React from 'react';
import { Link } from 'react-router-dom';

const Navbar = () => {
    return (
        <nav style={{
            backgroundColor: 'var(--surface-color)',
            boxShadow: 'var(--shadow-sm)',
            padding: '1rem 0',
            position: 'sticky',
            top: 0,
            zIndex: 100
        }}>
            <div className="container" style={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between'
            }}>
                {/* Logo */}
                <Link to="/" style={{
                    fontSize: '1.5rem',
                    fontWeight: 'bold',
                    color: 'var(--primary)',
                    display: 'flex',
                    alignItems: 'center',
                    gap: '0.5rem'
                }}>
                    <span>❤️</span> Lifeline
                </Link>

                {/* Navigation */}
                <div style={{ display: 'flex', gap: '2rem', alignItems: 'center' }}>
                    <Link to="/" style={{ fontWeight: 500, color: 'var(--text-secondary)' }}>Home</Link>
                    <Link to="/banks" style={{ fontWeight: 500, color: 'var(--text-secondary)' }}>Blood Banks</Link>
                    <a href="#" style={{ fontWeight: 500, color: 'var(--text-secondary)' }}>About</a>
                </div>

                {/* Auth Buttons */}
                <div style={{ display: 'flex', gap: '1rem' }}>
                    <Link to="/login" className="btn btn-secondary">
                        Login
                    </Link>
                    <Link to="/register" className="btn btn-primary">
                        Register
                    </Link>
                </div>
            </div>
        </nav>
    );
};

export default Navbar;
