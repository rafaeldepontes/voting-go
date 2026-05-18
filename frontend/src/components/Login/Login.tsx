import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { LogIn, UserPlus, AlertCircle, Loader2 } from 'lucide-react';
import styles from './Login.module.css';

interface LoginProps {
  onLoginSuccess: (token: string) => void;
}

const API_URL = import.meta.env.VITE_API_URL;

export const Login = ({ onLoginSuccess }: LoginProps) => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email || !password) {
      setError('Please fill in all fields');
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const response = await fetch(`${API_URL}/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || 'Login failed');
      }

      const data = await response.json();
      onLoginSuccess(data.token);
      navigate('/');
    } catch (err: any) {
      setError(err.message || 'Could not log in. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.loginPage}>
      <div className={styles.loginCard}>
        <div className={styles.cardHeader}>
          <div className={styles.iconWrapper}>
            <LogIn size={24} className={styles.icon} />
          </div>
          <h1>Welcome Back</h1>
          <p>Sign in to manage your polls and votes</p>
        </div>

        <form onSubmit={handleSubmit} className={styles.loginForm}>
          <div className={styles.formGroup}>
            <label htmlFor="email">Email Address</label>
            <input
              id="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="name@example.com"
              disabled={loading}
              required
            />
          </div>

          <div className={styles.formGroup}>
            <label htmlFor="password">Password</label>
            <input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="••••••••"
              disabled={loading}
              required
            />
          </div>

          {error && (
            <div className={styles.errorBanner}>
              <AlertCircle size={18} />
              <span>{error}</span>
            </div>
          )}

          <button type="submit" className={styles.loginButton} disabled={loading}>
            {loading ? <Loader2 size={18} className={styles.spinner} /> : 'Sign In'}
          </button>
        </form>

        <div className={styles.cardFooter}>
          <p>Don't have an account?</p>
          <Link to="/register" className={styles.linkButton}>
            <UserPlus size={16} />
            Create professional account
          </Link>
        </div>
      </div>
    </div>
  );
};
