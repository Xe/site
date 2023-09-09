export type Author = {
  name: string;
  handle: string;
  image?: string;
  url?: string;
  sameAs?: string[];
  jobTitle?: string;
  inSystem: boolean;
  pronouns: PronounSet;
};

export type Character = {
  name: string;
  stickerName: string;
  defaultPose: string;
  description: string;
  pronouns: PronounSet;
  stickers: string[];
};

export type Company = {
  name: string;
  url?: string;
  tagline: string;
  location: Location;
  defunct?: boolean;
};

export type Job = {
  company: Company;
  title: string;
  contract: boolean;
  startDate: string;
  endDate?: string;
  daysWorked?: number;
  daysBetween?: number;
  salary: Salary;
  leaveReason?: string;
  locations: Location[];
  highlights: string[];
  hideFromResume?: boolean;
};

export type Link = {
  url: string;
  title: string;
  description?: string;
};

export type Location = {
  city: string;
  stateOrProvince: string;
  country: string;
  remote: boolean;
};

export type NagMessage = {
  name: string;
  mood: string;
  message: string;
};

export type Person = {
  name: string;
  tags: string[];
  links: Link[];
};

export type PronounSet = {
  nominative: string;
  accusative: string;
  possessiveDeterminer: string;
  possessive: string;
  reflexive: string;
  singular: boolean;
};

export type Resume = {
  name: string;
  tagline: string;
  location: Location;
  buzzwords: string[];
  jobs: Job[];
  notablePublications: Link[];
};

export type Salary = {
  amount: number;
  currency: string;
  per: string;
  stock?: Stock;
};

export type Stock = {
  kind: "Grant" | "Options";
  amount: number;
  liquid: boolean;
  vestingYears: number;
  cliffYears: number;
};

export type StreamVOD = {
  title: string;
  slug: string;
  date: string;
  description: string;
  cdnPath: string;
  tags: string[];
};