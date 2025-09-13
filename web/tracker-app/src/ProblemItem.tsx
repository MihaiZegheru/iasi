import React from 'react';

export interface ProblemItemProps {
  name: string;
  url: string;
  time: string;
  solved: boolean;
  onToggle: () => void;
}

const ProblemItem: React.FC<ProblemItemProps> = ({ name, url, time, solved, onToggle }) => (
  <li className="problem-item">
    <input type="checkbox" checked={solved} onChange={onToggle} />
    <a
      href={url}
      target="_blank"
      rel="noopener noreferrer"
      className={solved ? 'solved url-link' : 'url-link'}
      style={{ textDecoration: 'none', color: solved ? '#888' : '#0074d9', fontWeight: 500 }}
    >
      {name}
    </a>
    <span style={{ marginLeft: 8, color: '#888' }}>({time})</span>
  </li>
);

export default ProblemItem;
