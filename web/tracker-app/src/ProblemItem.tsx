import React from 'react';
import { Link } from 'react-router-dom';

export interface ProblemItemProps {
  name: string;
  url: string;
  time: string;
  id: string;
  solved: boolean;
  onToggle: () => void;
}

const ProblemItem: React.FC<ProblemItemProps> = ({ name, url, time, id, solved, onToggle }) => (
  <li
    className="problem-item"
    style={{ cursor: 'pointer', position: 'relative' }}
    onClick={e => {
      // Prevent navigation if clicking on name or View Code
      const target = e.target as HTMLElement;
      if (target.closest('.problem-name') || target.closest('.view-code-link')) return;
      window.location.href = `/problem/${id}`;
    }}
  >
    <input type="checkbox" checked={solved} onChange={onToggle} />
    <a
      href={url}
      target="_blank"
      rel="noopener noreferrer"
      className={solved ? 'solved url-link problem-name' : 'url-link problem-name'}
      style={{ textDecoration: 'none', color: solved ? '#888' : '#0074d9', fontWeight: 500 }}
      onClick={e => e.stopPropagation()}
    >
      {name}
    </a>
    <span style={{ marginLeft: 8, color: '#888' }}>({time})</span>
    <a
      href={`https://www.infoarena.ro/job_detail/${id}?action=view-source`}
      target="_blank"
      rel="noopener noreferrer"
      className="view-code-link"
      style={{ marginLeft: 12, fontSize: '0.95em', color: '#555', textDecoration: 'underline' }}
      onClick={e => e.stopPropagation()}
    >
      View Code
    </a>
  </li>
);

export default ProblemItem;
