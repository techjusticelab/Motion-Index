import { writable, derived } from 'svelte/store';

export interface ValidationRule {
  required?: boolean;
  minLength?: number;
  maxLength?: number;
  pattern?: RegExp;
  custom?: (value: any) => string | null;
}

export interface FormState<T = Record<string, any>> {
  values: T;
  errors: Record<string, string>;
  touched: Record<string, boolean>;
  isSubmitting: boolean;
  isValid: boolean;
}

/**
 * Composable hook for form management with validation
 * Provides form state, validation, and submission handling
 */
export function useForm<T extends Record<string, any>>(
  initialValues: T,
  validationRules: Record<keyof T, ValidationRule> = {}
) {
  const formStore = writable<FormState<T>>({
    values: { ...initialValues },
    errors: {},
    touched: {},
    isSubmitting: false,
    isValid: true
  });

  // Derived store for easy access to just the values
  const values = derived(formStore, $form => $form.values);
  
  // Derived store for validation state
  const isValid = derived(formStore, $form => $form.isValid);

  /**
   * Validate a single field
   */
  function validateField(name: keyof T, value: any): string | null {
    const rules = validationRules[name];
    if (!rules) return null;

    // Required validation
    if (rules.required && (!value || value.toString().trim() === '')) {
      return `${String(name)} is required`;
    }

    // Skip other validations if value is empty and not required
    if (!value || value.toString().trim() === '') {
      return null;
    }

    // Min length validation
    if (rules.minLength && value.toString().length < rules.minLength) {
      return `${String(name)} must be at least ${rules.minLength} characters`;
    }

    // Max length validation
    if (rules.maxLength && value.toString().length > rules.maxLength) {
      return `${String(name)} must be no more than ${rules.maxLength} characters`;
    }

    // Pattern validation
    if (rules.pattern && !rules.pattern.test(value.toString())) {
      return `${String(name)} format is invalid`;
    }

    // Custom validation
    if (rules.custom) {
      return rules.custom(value);
    }

    return null;
  }

  /**
   * Validate all fields
   */
  function validateAll(): Record<string, string> {
    const errors: Record<string, string> = {};
    
    formStore.update(form => {
      Object.keys(form.values).forEach(key => {
        const error = validateField(key as keyof T, form.values[key as keyof T]);
        if (error) {
          errors[key] = error;
        }
      });

      return {
        ...form,
        errors,
        isValid: Object.keys(errors).length === 0
      };
    });

    return errors;
  }

  return {
    subscribe: formStore.subscribe,
    values,
    isValid,

    // Set field value and validate
    setField: (name: keyof T, value: any) => {
      formStore.update(form => {
        const newValues = { ...form.values, [name]: value };
        const error = validateField(name, value);
        const newErrors = { ...form.errors };
        
        if (error) {
          newErrors[String(name)] = error;
        } else {
          delete newErrors[String(name)];
        }

        return {
          ...form,
          values: newValues,
          errors: newErrors,
          touched: { ...form.touched, [name]: true },
          isValid: Object.keys(newErrors).length === 0
        };
      });
    },

    // Set multiple field values
    setFields: (fields: Partial<T>) => {
      formStore.update(form => {
        const newValues = { ...form.values, ...fields };
        const newErrors = { ...form.errors };
        const newTouched = { ...form.touched };

        Object.keys(fields).forEach(key => {
          const error = validateField(key as keyof T, fields[key as keyof T]);
          newTouched[key] = true;
          
          if (error) {
            newErrors[key] = error;
          } else {
            delete newErrors[key];
          }
        });

        return {
          ...form,
          values: newValues,
          errors: newErrors,
          touched: newTouched,
          isValid: Object.keys(newErrors).length === 0
        };
      });
    },

    // Mark field as touched
    setTouched: (name: keyof T, touched = true) => {
      formStore.update(form => ({
        ...form,
        touched: { ...form.touched, [name]: touched }
      }));
    },

    // Validate all fields
    validate: validateAll,

    // Reset form to initial values
    reset: () => {
      formStore.set({
        values: { ...initialValues },
        errors: {},
        touched: {},
        isSubmitting: false,
        isValid: true
      });
    },

    // Submit handler with validation
    handleSubmit: async (onSubmit: (values: T) => Promise<void> | void) => {
      const errors = validateAll();
      
      if (Object.keys(errors).length > 0) {
        // Mark all fields as touched to show errors
        formStore.update(form => ({
          ...form,
          touched: Object.keys(form.values).reduce(
            (acc, key) => ({ ...acc, [key]: true }), 
            {}
          )
        }));
        return;
      }

      formStore.update(form => ({ ...form, isSubmitting: true }));

      try {
        await onSubmit(formStore.values);
      } finally {
        formStore.update(form => ({ ...form, isSubmitting: false }));
      }
    }
  };
}