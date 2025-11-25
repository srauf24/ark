# ARK Design System

**Version**: 2.0 (Minimal Luxury)
**Last Updated**: November 25, 2025
**Theme**: S-Tier SaaS Dark Mode (True Black)

---

## Design Philosophy

ARK's design system embodies **Minimal Luxury**—a high-end, S-tier SaaS aesthetic inspired by Linear, Vercel, and Stripe. It prioritizes **meticulous craft, clarity, and restraint**.

### Core Principles

1.  **Deep Space Black**: True black (`#000000`) backgrounds for premium OLED aesthetics and maximum contrast.
2.  **Crisp White Text**: Pure white text for razor-sharp readability.
3.  **Selective Accents**: Electric blue is used *only* for interaction and focus, never for decoration.
4.  **No Clutter**: Flat surfaces, subtle borders, and zero unnecessary effects.
5.  **Purposeful Interaction**: Glows and lifts occur only on user interaction (hover/focus).

---

## Color System

All colors use the **OKLCH color space** for perceptual uniformity.

### Background Layers (True Black)

| Variable | Value | Hex Approx | Usage |
|----------|-------|------------|-------|
| `--background` | `oklch(0 0 0)` | `#000000` | Page background (True Black) |
| `--surface-0` | `oklch(0.06 0 0)` | `#0f0f0f` | Barely elevated (Inputs) |
| `--surface-1` | `oklch(0.10 0 0)` | `#1a1a1a` | Primary Cards/Panels |
| `--surface-2` | `oklch(0.14 0 0)` | `#242424` | Hover states |
| `--surface-3` | `oklch(0.18 0 0)` | `#2e2e2e` | Active/Selected |

### Text & Foreground

| Variable | Value | Hex Approx | Usage |
|----------|-------|------------|-------|
| `--foreground` | `oklch(1 0 0)` | `#ffffff` | Primary text |
| `--foreground-secondary` | `oklch(0.70 0 0)` | `#b3b3b3` | Secondary text |
| `--foreground-muted` | `oklch(0.50 0 0)` | `#808080` | Muted text |
| `--foreground-subtle` | `oklch(0.35 0 0)` | `#595959` | Subtle details |

### Electric Blue Accent

| Variable | Value | Usage |
|----------|-------|-------|
| `--primary` | `oklch(0.65 0.28 250)` | Primary Actions, Focus Rings |
| `--primary-foreground` | `oklch(0 0 0)` | Text on primary buttons |

**Usage Rule**: Use sparingly. Only for the primary call-to-action on a page or active states.

### Semantic Colors (Refined)

| Color | Value | Usage |
|-------|-------|-------|
| Success | `oklch(0.75 0.20 145)` | Valid states, success badges |
| Warning | `oklch(0.75 0.20 85)` | Warnings |
| Error | `oklch(0.65 0.22 25)` | Errors, destructive actions |
| Info | `oklch(0.65 0.28 250)` | Information |

---

## Typography

**Font Family**: `Inter`, system-ui, sans-serif.
**Monospace**: `JetBrains Mono`, monospace.

### Type Scale

| Class | Size | Line Height | Usage |
|-------|------|-------------|-------|
| `text-xs` | 12px | 1.5 | Metadata, Captions |
| `text-sm` | 14px | 1.5 | Body Small, UI Elements |
| `text-base` | 16px | 1.5 | Body Medium (Default) |
| `text-lg` | 18px | 1.5 | Body Large |
| `text-xl` | 20px | 1.4 | Section Headers |
| `text-2xl` | 24px | 1.3 | Page Headers |

### Weights

-   **Regular (400)**: Body text
-   **Medium (500)**: Buttons, Labels, Navigation
-   **Semibold (600)**: Headings

---

## Component Patterns

### Buttons

**Primary**:
-   Bg: `--primary` (Electric Blue)
-   Text: Black
-   Hover: Slight lift (`-1px`), subtle glow.
-   Radius: `4px`

**Secondary**:
-   Bg: Transparent
-   Border: `--border-default`
-   Hover: `--surface-1` background.

### Cards

-   Bg: `--surface-1` (`#1a1a1a`)
-   Border: `--border-subtle`
-   Radius: `6px`
-   Shadow: None (Flat)
-   Hover: Border becomes `--border-default`, slight lift.

### Inputs

-   Bg: `--surface-0` (`#0f0f0f`)
-   Border: `--border-default`
-   Focus: `--primary` border + 2px glow ring.
-   Radius: `4px`

### Badges

-   Bg: `--surface-2`
-   Text: `--foreground-secondary`
-   Radius: `2px`
-   Text: Uppercase, Medium weight, tracking wide.

---

## Spacing & Layout

**Base Unit**: 8px.

| Class | Value |
|-------|-------|
| `gap-1` | 4px |
| `gap-2` | 8px |
| `gap-3` | 12px |
| `gap-4` | 16px |
| `gap-6` | 24px |
| `gap-8` | 32px |

---

## Interaction Design

-   **Micro-interactions**: Fast (150ms).
-   **Hover**: Subtle lift (`transform: translateY(-1px)`).
-   **Focus**: Clear, accessible focus rings on all interactive elements.
-   **No Distractions**: No constant pulsing or animations on static elements.

---

## Accessibility

-   **Contrast**: Meets WCAG AA+ standards.
-   **Keyboard**: Full keyboard navigability.
-   **Focus**: Visible focus indicators.
