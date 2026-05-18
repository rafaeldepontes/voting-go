import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Vote, List, PlusCircle, LogIn, UserPlus, LogOut } from 'lucide-react';
import styles from './Navbar.module.css';

interface NavbarProps {
  token: string | null;
  onLogout: () => void;
}

export const Navbar = ({ token, onLogout }: NavbarProps) => {
  const navigate = useNavigate();

  return (
    <nav className={styles.navbar}>
      <div className={`container ${styles.navContainer}`}>
        <Link to="/" className={styles.logo}>
          <Vote size={24} className={styles.logoIcon} />
          <span>Voting Go</span>
        </Link>

        <div className={styles.navLinks}>
          <Link to="/" className={styles.navLink}>
            <List size={18} />
            <span>Polls</span>
          </Link>
          <Link to="/create" className={styles.navLink}>
            <PlusCircle size={18} />
            <span>Create Poll</span>
          </Link>
        </div>

        <div className={styles.authLinks}>
          {token ? (
            <button onClick={onLogout} className={styles.logoutButton}>
              <LogOut size={18} />
              <span>Logout</span>
            </button>
          ) : (
            <>
              <Link to="/login" className={styles.loginLink}>
                <LogIn size={18} />
                <span>Login</span>
              </Link>
              <Link to="/register" className={styles.registerLink}>
                <UserPlus size={18} />
                <span>Register</span>
              </Link>
            </>
          )}
        </div>
      </div>
    </nav>
  );
};
