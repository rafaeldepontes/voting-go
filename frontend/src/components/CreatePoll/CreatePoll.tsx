import { useState } from 'react';
import { Plus, Trash2, Send } from 'lucide-react';
import styles from './CreatePoll.module.css';

interface CreatePollProps {
  onPollCreated: () => void;
}

const API_URL = import.meta.env.VITE_API_URL;

export const CreatePoll = ({ onPollCreated }: CreatePollProps) => {
  const [name, setName] = useState('');
  const [options, setOptions] = useState(['', '']);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleAddOption = () => {
    setOptions([...options, '']);
  };

  const handleRemoveOption = (index: number) => {
    if (options.length <= 2) return;
    const newOptions = options.filter((_, i) => i !== index);
    setOptions(newOptions);
  };

  const handleOptionChange = (index: number, value: string) => {
    const newOptions = [...options];
    newOptions[index] = value;
    setOptions(newOptions);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim() || options.some(opt => !opt.trim())) {
      setError('Please fill in all fields');
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const response = await fetch(`${API_URL}/polls`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          name: name.trim(),
          options: options.map(opt => opt.trim()),
        }),
      });

      if (!response.ok) throw new Error('Failed to create poll');
      
      onPollCreated();
    } catch (err) {
      setError('Could not create poll. Is the backend running?');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.createPoll}>
      <h2>Create a New Poll</h2>
      <form onSubmit={handleSubmit} className={styles.pollForm}>
        <div className={styles.formGroup}>
          <label htmlFor="poll-name">Poll Question</label>
          <input
            id="poll-name"
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="e.g., What is your favorite programming language?"
            disabled={loading}
          />
        </div>

        <div className={styles.formGroup}>
          <label>Options</label>
          {options.map((option, index) => (
            <div key={index} className={styles.optionInputGroup}>
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
                className={styles.deleteButton}
                disabled={loading || options.length <= 2}
              >
                <Trash2 size={18} />
              </button>
            </div>
          ))}
          <button
            type="button"
            onClick={handleAddOption}
            className={styles.addOptionButton}
            disabled={loading}
          >
            <Plus size={18} /> Add Option
          </button>
        </div>

        {error && <p className={styles.error}>{error}</p>}

        <button type="submit" className={styles.submitButton} disabled={loading}>
          {loading ? 'Creating...' : (
            <>
              <Send size={18} /> Create Poll
            </>
          )}
        </button>
      </form>
    </div>
  );
};
