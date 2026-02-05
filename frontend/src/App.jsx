import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Navbar from './components/Navbar';
import Home from './pages/Home';

// Placeholders for now
const Login = () => <div className="container" style={{ marginTop: '2rem' }}><h1>Login Page</h1></div>;
const Register = () => <div className="container" style={{ marginTop: '2rem' }}><h1>Register Page</h1></div>;
const Banks = () => <div className="container" style={{ marginTop: '2rem' }}><h1>Blood Banks</h1></div>;

function App() {
  return (
    <Router>
      <div style={{ minHeight: '100vh', display: 'flex', flexDirection: 'column' }}>
        <Navbar />
        <main style={{ flex: 1 }}>
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />
            <Route path="/banks" element={<Banks />} />
          </Routes>
        </main>

        {/* Simple Footer */}
        <footer style={{
          backgroundColor: 'white',
          padding: '2rem 0',
          textAlign: 'center',
          color: 'var(--text-secondary)',
          borderTop: '1px solid #E5E7EB'
        }}>
          <p>© 2026 Lifeline. Connecting donors, saving lives.</p>
        </footer>
      </div>
    </Router>
  );
}

export default App;
