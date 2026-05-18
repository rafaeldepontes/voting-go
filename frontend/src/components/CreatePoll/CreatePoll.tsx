import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { PlusCircle, Trash2, AlertCircle, Loader2, ArrowLeft, Send } from 'lucide-react';
import styles from './CreatePoll.module.css';

interface CreatePollProps {
  token: string | null;
  onAuthError: () => void;
  onPollCreated: () => void;
}

const API_URL = import.meta.env.VITE_API_URL;

const DURATION_OPTIONS = [
  { label: '1 Hour', value: 60 * 60 * 1000 * 1000 * 1000 },
  { label: '6 Hours', value: 6 * 60 * 60 * 1000 * 1000 * 1000 },
  { label: '24 Hours', value: 24 * 60 * 60 * 1000 * 1000 * 1000 },
  { label: '7 Days', value: 7 * 24 * 60 * 60 * 1000 * 1000 * 1000 },
  { label: '30 Days', value: 30 * 24 * 60 * 60 * 1000 * 1000 * 1000 },
  { label: 'Permanent', value: 0 },
];

export const CreatePoll = ({ token, onAuthError }: CreatePollProps) => {
  const [name, setName] = useState('');
  const [options, setOptions] = useState(['', '']);
  const [duration, setDuration] = useState(DURATION_OPTIONS[2].value); // Default 24h
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  const handleAddOption = () => {
    setOptions([...options, '']);
  };

  const handleRemoveOption = (index: number) => {
    if (options.length > 2) {
      const newOptions = options.filter((_, i) => i !== index);
      setOptions(newOptions);
    }
  };

  const handleOptionChange = (index: number, value: string) => {
    const newOptions = [...options];
    newOptions[index] = value;
    setOptions(newOptions);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) {
      setError('Poll question is required');
      return;
    }

    const filteredOptions = options.filter(opt => opt.trim() !== '');
    if (filteredOptions.length < 2) {
      setError('At least two options are required');
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const headers: Record<string, string> = { 'Content-Type': 'application/json' };
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }

      const response = await fetch(`${API_URL}/polls`, {
        method: 'POST',
        headers,
        body: JSON.stringify({
          name: name.trim(),
          options: filteredOptions.map(opt => opt.trim()),
          duration: duration,
        }),
      });

      if (response.status === 401 || response.status === 403) {
        if (token) onAuthError();
      }

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || 'Failed to create poll');
      }

      const data = await response.json();
      navigate(`/poll/${data.id}`);
    } catch (err: any) {
      setError(err.message || 'Could not create poll. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.createPollPage}>
      <button onClick={() => navigate(-1)} className={styles.backButton}>
        <ArrowLeft size={18} />
        <span>Back to Polls</span>
      </button>

      <div className={styles.createCard}>
        <div className={styles.cardHeader}>
          <h1>Create New Poll</h1>
          <p>Gather opinions from the community in real-time</p>
        </div>

        <form onSubmit={handleSubmit} className={styles.createForm}>
          <div className={styles.formSection}>
            <label htmlFor="poll-name">Poll Question</label>
            <input
              id="poll-name"
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="What would you like to ask?"
              disabled={loading}
              className={styles.questionInput}
            />
          </div>

          <div className={styles.formSection}>
            <label>Options</label>
            <div className={styles.optionsList}>
              {options.map((option, index) => (
                <div key={index} className={styles.optionItem}>
                  <input
                    type="text"
                    value={option}
                    onChange={(e) => handleOptionChange(index, e.target.value)}
                    placeholder={`Option ${index + 1}`}
                    disabled={loading}
                  />
                  <button
                    type="button"
                    onClick={() => handleRemoveOption(index)}
                    className={styles.removeOption}
                    disabled={loading || options.length <= 2}
                    title="Remove option"
                  >
                    <Trash2 size={18} />
                  </button>
                </div>
              ))}
            </div>
            <button
              type="button"
              onClick={handleAddOption}
              className={styles.addOptionButton}
              disabled={loading}
            >
              <PlusCircle size={18} />
              <span>Add another option</span>
            </button>
          </div>

          <div className={styles.formSection}>
            <label htmlFor="poll-duration">Poll Duration</label>
            <select
              id="poll-duration"
              value={duration}
              onChange={(e) => setDuration(Number(e.target.value))}
              disabled={loading}
              className={styles.durationSelect}
            >
              {DURATION_OPTIONS.map((opt) => (
                <option key={opt.value} value={opt.value}>
                  {opt.label}
                </option>
              ))}
            </select>
          </div>

          {error && (
            <div className={styles.errorBanner}>
              <AlertCircle size={18} />
              <span>{error}</span>
            </div>
          )}

          {!token && (
            <div className={styles.anonymousNote}>
              <p>You are creating this poll anonymously (identified by IP).</p>
            </div>
          )}

          <button type="submit" className={styles.submitButton} disabled={loading}>
            {loading ? <Loader2 size={18} className={styles.spinner} /> : (
              <>
                <Send size={18} />
                <span>Launch Poll</span>
              </>
            )}
          </button>
        </form>
      </div>
    </div>
  );
};
