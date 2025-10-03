import { writable } from 'svelte/store';

export interface FileUploadState {
  selectedFile: File | null;
  isDragging: boolean;
  isUploading: boolean;
  uploadStatus: string;
  uploadProgress: number;
}

export interface FileUploadOptions {
  allowedTypes?: string[];
  maxSizeBytes?: number;
  multiple?: boolean;
}

/**
 * Composable hook for file upload functionality
 * Handles file selection, drag & drop, validation, and upload state
 */
export function useFileUpload(options: FileUploadOptions = {}) {
  const {
    allowedTypes = ['.pdf', '.docx', '.txt'],
    maxSizeBytes = 50 * 1024 * 1024, // 50MB default
    multiple = false
  } = options;

  const state = writable<FileUploadState>({
    selectedFile: null,
    isDragging: false,
    isUploading: false,
    uploadStatus: '',
    uploadProgress: 0
  });

  /**
   * Validate file type and size
   */
  function validateFile(file: File): { valid: boolean; error?: string } {
    // Check file size
    if (file.size > maxSizeBytes) {
      return {
        valid: false,
        error: `File size must be less than ${(maxSizeBytes / 1024 / 1024).toFixed(1)}MB`
      };
    }

    // Check file type
    const fileExtension = '.' + file.name.split('.').pop()?.toLowerCase();
    if (!allowedTypes.some(type => type.toLowerCase() === fileExtension)) {
      return {
        valid: false,
        error: `File type not supported. Allowed types: ${allowedTypes.join(', ')}`
      };
    }

    return { valid: true };
  }

  /**
   * Handle file selection from input
   */
  function handleFileChange(event: Event) {
    const target = event.target as HTMLInputElement;
    const files = target.files;
    
    if (files && files.length > 0) {
      const file = files[0];
      const validation = validateFile(file);
      
      if (validation.valid) {
        state.update(s => ({ 
          ...s, 
          selectedFile: file, 
          uploadStatus: '' 
        }));
      } else {
        state.update(s => ({ 
          ...s, 
          selectedFile: null, 
          uploadStatus: validation.error || 'Invalid file' 
        }));
      }
    } else {
      state.update(s => ({ 
        ...s, 
        selectedFile: null, 
        uploadStatus: '' 
      }));
    }
  }

  /**
   * Handle drag over event
   */
  function handleDragOver(event: DragEvent) {
    event.preventDefault();
    state.update(s => ({ ...s, isDragging: true }));
  }

  /**
   * Handle drag leave event
   */
  function handleDragLeave(event: DragEvent) {
    event.preventDefault();
    state.update(s => ({ ...s, isDragging: false }));
  }

  /**
   * Handle file drop
   */
  function handleDrop(event: DragEvent) {
    event.preventDefault();
    state.update(s => ({ ...s, isDragging: false }));

    const files = event.dataTransfer?.files;
    if (files && files.length > 0) {
      const file = files[0];
      const validation = validateFile(file);
      
      if (validation.valid) {
        state.update(s => ({ 
          ...s, 
          selectedFile: file, 
          uploadStatus: '' 
        }));
      } else {
        state.update(s => ({ 
          ...s, 
          selectedFile: null, 
          uploadStatus: validation.error || 'Invalid file' 
        }));
      }
    }
  }

  /**
   * Clear selected file
   */
  function clearFile() {
    state.update(s => ({ 
      ...s, 
      selectedFile: null, 
      uploadStatus: '',
      uploadProgress: 0
    }));
  }

  /**
   * Set upload status
   */
  function setUploadStatus(status: string, isError = false) {
    state.update(s => ({ ...s, uploadStatus: status }));
  }

  /**
   * Set upload progress
   */
  function setUploadProgress(progress: number) {
    state.update(s => ({ ...s, uploadProgress: Math.max(0, Math.min(100, progress)) }));
  }

  /**
   * Start upload process
   */
  function startUpload() {
    state.update(s => ({ ...s, isUploading: true, uploadProgress: 0 }));
  }

  /**
   * Complete upload process
   */
  function completeUpload(success: boolean, message?: string) {
    state.update(s => ({ 
      ...s, 
      isUploading: false, 
      uploadProgress: success ? 100 : 0,
      uploadStatus: message || (success ? 'Upload completed' : 'Upload failed')
    }));
  }

  /**
   * Reset all state
   */
  function reset() {
    state.set({
      selectedFile: null,
      isDragging: false,
      isUploading: false,
      uploadStatus: '',
      uploadProgress: 0
    });
  }

  return {
    subscribe: state.subscribe,
    
    // Event handlers
    handleFileChange,
    handleDragOver,
    handleDragLeave,
    handleDrop,
    
    // Actions
    clearFile,
    setUploadStatus,
    setUploadProgress,
    startUpload,
    completeUpload,
    reset,
    
    // Utilities
    validateFile,
    
    // Configuration
    allowedTypes,
    maxSizeBytes
  };
}