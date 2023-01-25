import { useState, useEffect } from "react";
import axios from 'axios';
import {useParams} from "react-router-dom";

function StockOverview() {
  // We can use the `useParams` hook here to access
  // the dynamic pieces of the URL.
  let { id } = useParams();

  return (
    <div>
      <h3>ID: {id}</h3>
    </div>
  );
}

export default StockOverview