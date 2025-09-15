import React, { useState } from 'react';

interface AccordionBoxProps {
  title: string;
  children: React.ReactNode;
  defaultOpen?: boolean;
}

const AccordionBox: React.FC<AccordionBoxProps> = ({ title, children, defaultOpen = false }) => {
  const [open, setOpen] = useState(defaultOpen);
  return (
    <div className={`accordion-box${open ? ' open' : ''}`}>
      <div
        className="accordion-title"
        onClick={() => setOpen(o => !o)}
        style={{ userSelect: 'none', cursor: 'pointer', textAlign: 'center', position: 'relative' }}
      >
        <span>{title}</span>
        <span style={{ fontWeight: 400, position: 'absolute', right: 18, top: '50%', transform: 'translateY(-50%)' }}>{open ? '▲' : '▼'}</span>
      </div>
      {open && (
        <div className="accordion-content">
          {children}
        </div>
      )}
    </div>
  );
};

export default AccordionBox;
