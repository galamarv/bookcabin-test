import React, { useState } from 'react';
import './App.css';

/**
 * App component for generating seat vouchers for crew members. 
 * It allows users to input crew details and flight information,
 * checks if vouchers have already been generated for the specified flight,
 * and generates seat assignments if not. 
 * It handles form submission, input changes, and displays results or errors.
 */
function App() {
  // State to manage form data for crew and flight details.
  // formData holds the input values for crew name, ID, flight number, date,
  // and aircraft type. It initializes with default values.
  const [formData, setFormData] = useState({
    name: '',
    id: '',
    flightNumber: '',
    date: '',
    aircraft: 'ATR',
  });

  // State to manage form data, assigned seats, error messages, and loading state.
  const [assignedSeats, setAssignedSeats] = useState([]);
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  // handleInputChange updates the formData state when the user types in the input fields.
  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });
  };

  // handleSubmit is called when the form is submitted.
  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsLoading(true);
    setError('');
    setAssignedSeats([]);

    if (!formData.name || !formData.id || !formData.flightNumber || !formData.date) {
      setError('All fields are required.');
      setIsLoading(false);
      return;
    }

    try {
      const checkRes = await fetch('/api/check', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          flightNumber: formData.flightNumber,
          date: formData.date,
        }),
      });

      const checkData = await checkRes.json();

      if (checkData.exists) {
        setError('Vouchers have already been generated for this flight on this date.');
        setIsLoading(false);
        return;
      }

      const generateRes = await fetch('/api/generate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(formData),
      });

      const generateData = await generateRes.json();

      if (generateData.success) {
        setAssignedSeats(generateData.seats);
      } else {
        setError(generateData.error || 'An unknown error occurred.');
      }

    } catch (err) {
      setError('Failed to connect to the server. Please ensure the backend is running.');
    } finally {
      setIsLoading(false);
    }
  };

  // The main render function of the App component.
  return (
    <div className="App">
      <div className="app-header">
        <h1>Voucher Seat Assignment</h1>
        <p>Enter flight and crew details to generate seat vouchers.</p>
      </div>

      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="name">Crew Name</label>
          <input type="text" id="name" name="name" value={formData.name} onChange={handleInputChange} />
        </div>
        <div className="form-group">
          <label htmlFor="id">Crew ID</label>
          <input type="text" id="id" name="id" value={formData.id} onChange={handleInputChange} />
        </div>
        <div className="form-group">
          <label htmlFor="flightNumber">Flight Number</label>
          <input type="text" id="flightNumber" name="flightNumber" value={formData.flightNumber} onChange={handleInputChange} />
        </div>
        <div className="form-group">
          <label htmlFor="date">Flight Date</label>
          <input type="date" id="date" name="date" value={formData.date} onChange={handleInputChange} />
        </div>
        <div className="form-group">
          <label htmlFor="aircraft">Aircraft Type</label>
          <select id="aircraft" name="aircraft" value={formData.aircraft} onChange={handleInputChange}>
            <option value="ATR">ATR</option>
            <option value="Airbus 320">Airbus 320</option>
            <option value="Boeing 737 Max">Boeing 737 Max</option>
          </select>
        </div>
        <button type="submit" className="btn-submit" disabled={isLoading}>
          {isLoading ? 'Generating...' : 'Generate Vouchers'}
        </button>
      </form>

      {error && <div className="message error">{error}</div>}

      {assignedSeats.length > 0 && (
        <div className="results">
          <h2>Assigned Voucher Seats</h2>
          <ul className="seat-list">
            {assignedSeats.map((seat) => (
              <li key={seat} className="seat-item">{seat}</li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
}

export default App;