import { createContext, useContext } from 'react';

export const ServicesContext = createContext();

const useServices = () => {
  const services = useContext(ServicesContext);

  return services;
};

export default useServices;
