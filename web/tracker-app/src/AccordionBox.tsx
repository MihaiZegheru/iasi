import React, { useState } from 'react';

interface AccordionBoxProps {
  title: string;
  children: React.ReactNode;
  defaultOpen?: boolean;
  locked?: boolean;
  onUnlock?: () => void;
}

const AccordionBox: React.FC<AccordionBoxProps> = ({ title, children, defaultOpen = false, locked = false, onUnlock }) => {
  const [open, setOpen] = useState(defaultOpen);
  return (
    <div className={`accordion-box${open ? ' open' : ''}`}>
      <div
        className="accordion-title"
        onClick={() => setOpen(o => !o)}
        style={{ userSelect: 'none', cursor: 'pointer', textAlign: 'center', position: 'relative' }}
      >
        <span>{title}</span>
        {locked ? (
          <span
            className="accordion-lock"
            title="Locked. Click to unlock."
            onClick={e => {
              e.stopPropagation();
              if (onUnlock) onUnlock();
            }}
            style={{ cursor: 'pointer', marginLeft: 8, position: 'absolute', left: 18, top: '50%', transform: 'translateY(-50%)' }}
          >
            <span role="img" aria-label="locked">🔒</span>
          </span>
        ) : null}
        <span style={{ fontWeight: 400, position: 'absolute', right: 18, top: '50%', transform: 'translateY(-50%)' }}>{open ? '▲' : '▼'}</span>
      </div>
      {open && (
        <div className="accordion-content">
          {locked ? (
            <div className="accordion-locked-message">
              <span style={{ opacity: 0.7 }}>
                This content is locked. Click the lock to unlock.
              </span>
            </div>
          ) : (
            children
          )}
        </div>
      )}
    </div>
  );
};

export default AccordionBox;
