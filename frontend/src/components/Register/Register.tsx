import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { UserPlus, LogIn, AlertCircle, CheckCircle, Loader2 } from 'lucide-react';
import styles from './Register.module.css';

interface RegisterProps {
  onRegisterSuccess: () => void;
}

const API_URL = import.meta.env.VITE_API_URL;

export const Register = ({ onRegisterSuccess }: RegisterProps) => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    if (!email || !password || !confirmPassword) {
      setError('Please fill in all fields');
      return;
    }

    if (password !== confirmPassword) {
      setError('Passwords do not match');
      return;
    }

    if (password.length < 5) {
        setError('Password must be at least 5 characters long');
        return;
    }

    setLoading(true);

    try {
      const response = await fetch(`${API_URL}/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || 'Registration failed');
      }

      setSuccess(true);
      onRegisterSuccess();
      setTimeout(() => {
        navigate('/login');
      }, 2000);
    } catch (err: any) {
      setError(err.message || 'Could not register. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  if (success) {
    return (
      <div className={styles.registerPage}>
        <div className={`${styles.registerCard} ${styles.successCard}`}>
          <div className={styles.successIconWrapper}>
            <CheckCircle size={48} className={styles.successIcon} />
          </div>
          <h2>Registration Successful!</h2>
          <p>Your account has been created. Redirecting to login...</p>
        </div>
      </div>
    );
  }

  return (
    <div className={styles.registerPage}>
      <div className={styles.registerCard}>
        <div className={styles.cardHeader}>
          <div className={styles.iconWrapper}>
            <UserPlus size={24} className={styles.icon} />
          </div>
          <h1>Create Account</h1>
          <p>Join Voting Go to start creating professional polls</p>
        </div>

        <form onSubmit={handleSubmit} className={styles.registerForm}>
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

          <div className={styles.formGroup}>
            <label htmlFor="confirmPassword">Confirm Password</label>
            <input
              id="confirmPassword"
              type="password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
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

          <button type="submit" className={styles.registerButton} disabled={loading}>
            {loading ? <Loader2 size={18} className={styles.spinner} /> : 'Create Account'}
          </button>
        </form>

        <div className={styles.cardFooter}>
          <p>Already have an account?</p>
          <Link to="/login" className={styles.linkButton}>
            <LogIn size={16} />
            Sign in instead
          </Link>
        </div>
      </div>
    </div>
  );
};
