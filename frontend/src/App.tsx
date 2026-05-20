import { useState, useEffect } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { Navbar } from './components/Layout/Navbar';
import { PollList } from './components/PollList/PollList';
import { CreatePoll } from './components/CreatePoll/CreatePoll';
import { ResultsView } from './components/ResultsView/ResultsView';
import { VoteView } from './components/VoteView/VoteView';
import { Login } from './components/Login/Login';
import { Register } from './components/Register/Register';

function App() {
  const [token, setToken] = useState<string | null>(localStorage.getItem('token'));

  useEffect(() => {
    if (token) {
      localStorage.setItem('token', token);
    } else {
      localStorage.removeItem('token');
    }
  }, [token]);

  const handleLogout = () => {
    setToken(null);
  };

  const handleAuthError = () => {
    // Idk why I did this... But now it became troublesome...
    // setToken(null);
  };

  return (
    <BrowserRouter>
      <Navbar token={token} onLogout={handleLogout} />

      <main className="content container">
        <Routes>
          <Route path="/" element={<PollList token={token} onAuthError={handleAuthError} />} />
          <Route path="/login" element={<Login onLoginSuccess={(t) => setToken(t)} />} />
          <Route path="/register" element={<Register onRegisterSuccess={() => { }} />} />
          <Route path="/create" element={<CreatePoll token={token} onAuthError={handleAuthError} onPollCreated={() => { }} />} />
          <Route path="/poll/:id" element={<VoteView token={token} onAuthError={handleAuthError} onVoteSuccess={() => { }} />} />
          <Route path="/results/:id" element={<ResultsView token={token} />} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </main>

      <footer style={{ padding: '2rem 0', textAlign: 'center', color: 'var(--text-muted)', fontSize: '0.875rem', borderTop: '1px solid var(--border)', marginTop: 'auto' }}>
        <p>&copy; 2026 Voting Go · Built by Rafael</p>
      </footer>
    </BrowserRouter>
  );
}

export default App;
